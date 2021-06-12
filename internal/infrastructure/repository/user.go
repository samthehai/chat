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

// FindFriends implemented base on graphql connections
// https://relay.dev/graphql/connections.htm
func (r *UserRepository) FindFriends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error) {
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
		query =
			"SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified, " +
				"TRUE AS has_previous_page, " +

				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) " +
				"  FROM ( " +
				"   SELECT * FROM users " +
				"   ORDER BY " + sortColumn + " ASC, id ASC LIMIT $1 + 1 " +
				"  ) as np " +
				" ) = $1 + 1 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_next_page " +

				"FROM users " +
				"ORDER BY " + sortColumn + " ASC, id ASC LIMIT $1"

		rows, err = r.db.QueryContext(ctx, query, first)
	} else {
		query =
			"SELECT id, name, picture_url, firebase_id, provider, email_address, email_verified, " +
				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) FROM users " +
				"   WHERE " + sortColumn + " <= (SELECT " + sortColumn + " FROM users WHERE id = $2 ) " +
				"   AND id != $2 " +
				"   AND id NOT IN " +
				"    (SELECT id FROM users " +
				"      WHERE " + sortColumn + " = (SELECT " + sortColumn + " FROM users WHERE id = $2 ) " +
				"      AND id >= $2 ) " +
				" ) > 0 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_previous_page, " +

				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) " +
				"  FROM ( " +
				"   SELECT * FROM users " +
				"   WHERE " + sortColumn + " >= (SELECT " + sortColumn + " FROM users WHERE id = $2 ) " +
				"   AND id != $2 " +
				"   ORDER BY " + sortColumn + " ASC, id ASC LIMIT $1 + 1 " +
				"  ) as np " +
				" ) = $1 + 1 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_next_page " +

				"FROM users " +
				"WHERE " + sortColumn + " >= (SELECT " + sortColumn + " FROM users WHERE id = $2 ) " +
				"AND id != $2 " +
				"AND id NOT IN ( " +
				" SELECT id FROM users " +
				" WHERE " + sortColumn + " = (SELECT " + sortColumn + " FROM users WHERE id = $2 ) " +
				" AND id <= $2 " +
				") " +
				"ORDER BY " + sortColumn + " ASC, id ASC LIMIT $1"

		rows, err = r.db.QueryContext(ctx, query, first, after)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userFriendEdges []*entity.UserFriendsEdge
		hasNextPage     bool
		hasPreviousPage bool
	)

	for rows.Next() {
		var edge struct {
			model.User
			HasNextPage     bool `json:"has_next_page"`
			HasPreviousPage bool `json:"has_previous_page"`
		}

		if err := rows.Scan(
			&edge.ID,
			&edge.Name,
			&edge.PictureUrl,
			&edge.FirebaseID,
			&edge.Provider,
			&edge.EmailAddress,
			&edge.EmailVerified,
			&edge.HasNextPage,
			&edge.HasPreviousPage,
		); err != nil {
			return nil, err
		}

		hasNextPage = edge.HasNextPage
		hasPreviousPage = edge.HasPreviousPage
		user := model.ConvertModelUser(&edge.User)
		userFriendEdges = append(userFriendEdges, &entity.UserFriendsEdge{
			Node:   *user,
			Cursor: user.ID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(userFriendEdges) == 0 {
		return &entity.UserFriendsConnection{}, nil
	}

	return &entity.UserFriendsConnection{
		Edges: userFriendEdges,
		PageInfo: entity.PageInfo{
			HasPreviousPage: hasPreviousPage,
			HasNextPage:     hasNextPage,
			StartCursor:     userFriendEdges[0].Cursor,
			EndCursor:       userFriendEdges[len(userFriendEdges)-1].Cursor,
		},
		TotalCount: int64(len(userFriendEdges)),
	}, nil
}
