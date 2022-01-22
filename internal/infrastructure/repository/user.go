package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/lib/pq"
	"github.com/samthehai/chat/internal/domain/entity"
	domainerrors "github.com/samthehai/chat/internal/domain/errors"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
	"github.com/samthehai/chat/internal/infrastructure/repository/model"
)

const usersKey = "users"

type UserRepository struct {
	cacher        external.Cacher
	authenticator external.Authenticator
	userChans     map[entity.ID]chan *entity.User
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
		userChans:     map[entity.ID]chan *entity.User{},
		mutex:         sync.Mutex{},
	}
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
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified FROM users ORDER BY distinct_id ASC LIMIT $1 OFFSET $2`, limit, offset)
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

func (r *UserRepository) FindUsers(ctx context.Context, userIDs []entity.ID) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified
		 FROM users
		 WHERE id = ANY($1)`,
		pq.Array(userIDs),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.FirebaseID, &user.Provider); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return model.ConvertModelUsers(users), nil
}

func (r *UserRepository) FindByFirebaseID(ctx context.Context, firebaseID string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified FROM users WHERE firebase_id=$1 LIMIT 1`, firebaseID)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.PictureUrl, &user.FirebaseID, &user.Provider, &user.EmailAddress, &user.EmailVerified)

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
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO users (name, picture_url, firebase_id, provider, email_address, email_verified) VALUES ($1,$2,$3,$4,$5,$6)`)
	if err != nil {
		return nil, fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, input.Name, input.PictureUrl, input.FirebaseID, input.Provider, input.EmailAddress, input.EmailVerified); err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}

	uj, _ := json.Marshal(&input)
	if err := r.cacher.LPush(usersKey, uj); err != nil {
		return nil, fmt.Errorf("lpush: %w", err)
	}

	createdUser, err := r.FindByFirebaseID(ctx, input.FirebaseID)
	if err != nil {
		return nil, fmt.Errorf("find by firebase id: %w", err)
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
		return nil, fmt.Errorf("get auth token from context")
	}

	return r.FindByFirebaseID(ctx, token.UserID)
}

func (r *UserRepository) GetAuthTokenFromContext(ctx context.Context) (*entity.AuthToken, error) {
	token, err := r.authenticator.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth token from context")
	}

	return token, nil
}

func (r *UserRepository) FindFriends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) ([]*entity.User, error) {
	if !entity.IsValidFriendsSortByType(string(sortBy)) {
		return nil, fmt.Errorf("invalid sortBy: %v", sortBy)
	}

	sortColumn := model.GetColumnNameByFriendsSortByType(sortBy)
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if after == 0 {
		query = fmt.Sprintf(
			` SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified
				FROM users
				ORDER BY %v ASC, id ASC LIMIT $1`,
			sortColumn,
		)
		rows, err = r.db.QueryContext(ctx, query, first)
	} else {
		query = fmt.Sprintf(
			` SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified
			  FROM users
			  WHERE %[1]v >= (SELECT %[1]v FROM users WHERE id = $2 )
					AND id != $2
					AND id NOT IN (
				  	SELECT id FROM users
				  	WHERE %[1]v = (SELECT %[1]v FROM users WHERE id = $2)
				  		AND id <= $2
					)
			  ORDER BY %[1]v ASC, id ASC LIMIT $1`,
			sortColumn,
		)

		rows, err = r.db.QueryContext(ctx, query, first, after)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.FirebaseID, &user.Provider); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return model.ConvertModelUsers(users), nil
}
