package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/samthehai/chat/internal/domain/entity"
	domainerrors "github.com/samthehai/chat/internal/domain/errors"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
	"github.com/samthehai/chat/internal/infrastructure/repository/model"
)

const usersKey = "users"

type UserRepository struct {
	cacher        external.Cacher
	authenticator external.Authenticator
	userChans     map[string]chan *entity.User
	mutex         sync.Mutex
	db            *sql.DB
}

func NewUserRepository(
	cacher external.Cacher,
	authenticator external.Authenticator,
	db *sql.DB,
) *UserRepository {
	return &UserRepository{
		cacher:        cacher,
		authenticator: authenticator,
		db:            db,
		userChans:     map[string]chan *entity.User{},
		mutex:         sync.Mutex{},
	}
}

func (r *UserRepository) Users(ctx context.Context) ([]*entity.User, error) {
	values, err := r.cacher.LRange(usersKey, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("[cacher] lrange: %w", err)
	}

	users := []*entity.User{}
	for _, uj := range values {
		u := &entity.User{}
		if err := json.Unmarshal([]byte(uj), &u); err != nil {
			return nil, fmt.Errorf("[Unmarshal] users: %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) UserJoined(ctx context.Context, input entity.User) (<-chan *entity.User, error) {
	users := make(chan *entity.User, 1)

	r.mutex.Lock()
	r.userChans[input.ID] = users
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()

		r.mutex.Lock()
		delete(r.userChans, input.ID)
		r.mutex.Unlock()
	}()

	return users, nil
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int64) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, firebase_id, provider FROM users ORDER BY distinct_id ASC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var user *model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.FirebaseID, &user.Provider); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return model.ConvertModelUsers(users), nil
}

func (r *UserRepository) Find(ctx context.Context, id string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, firebase_id, provider FROM users WHERE id=? LIMIT 1`, id)

	var user *model.User

	err := row.Scan(&user.ID, &user.Name, &user.FirebaseID, &user.Provider)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("%w: %s", domainerrors.ErrNotFound, err)
	case err != nil:
		return nil, err
	default:
		return model.ConvertModelUser(user), nil
	}
}

func (r *UserRepository) FindByFirebaseID(ctx context.Context, firebaseID string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, firebase_id, provider FROM users WHERE firebase_id=$1 LIMIT 1`, firebaseID)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.FirebaseID, &user.Provider)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("%w: %s", domainerrors.ErrNotFound, err)
	case err != nil:
		return nil, err
	default:
		return model.ConvertModelUser(&user), nil
	}
}

func (r *UserRepository) AddUser(ctx context.Context, input entity.User) (*entity.User, error) {
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO users (name, firebase_id, provider) VALUES ($1,$2,$3)`)
	if err != nil {
		return nil, fmt.Errorf("[DB] prepare context: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, input.Name, input.FirebaseID, input.Provider); err != nil {
		return nil, fmt.Errorf("[DB] exec context: %w", err)
	}

	uj, _ := json.Marshal(&input)
	if err := r.cacher.LPush(usersKey, uj); err != nil {
		return nil, fmt.Errorf("[cacher] lpush: %w", err)
	}

	createdUser, err := r.FindByFirebaseID(ctx, input.FirebaseID)
	if err != nil {
		return nil, fmt.Errorf("[DB] find by firebase id: %w", err)
	}

	r.mutex.Lock()
	for _, ch := range r.userChans {
		ch <- createdUser
	}
	r.mutex.Unlock()

	return createdUser, nil
}

func (r *UserRepository) GetUserFromContext(ctx context.Context) (*entity.User, error) {
	token, err := r.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Authenticator] get auth token from context")
	}

	return r.FindByFirebaseID(ctx, token.UserID)
}

func (r *UserRepository) GetAuthTokenFromContext(ctx context.Context) (*entity.AuthToken, error) {
	token, err := r.authenticator.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Authenticator] get auth token from context")
	}

	return token, nil
}
