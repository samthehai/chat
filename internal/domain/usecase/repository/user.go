package repository

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type UserRepository interface {
	AddUser(ctx context.Context, input entity.User) (*entity.User, error)
	FindByFirebaseID(ctx context.Context, firebaseID string) (*entity.User, error)
	GetUserFromContext(ctx context.Context) (*entity.User, error)
	GetAuthTokenFromContext(ctx context.Context) (*entity.AuthToken, error)
	UserJoined(ctx context.Context, user entity.User) (<-chan *entity.User, error)
	FindFriends(ctx context.Context, first int, after entity.ID,
		sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType,
	) (*entity.FriendsConnection, error)
	FindUsers(ctx context.Context, userIDs []entity.ID) ([]*entity.User, error)
	GetFriendIDsFromUserIDs(ctx context.Context,
		inputs []entity.FriendsQueryInput) (map[entity.ID]*entity.IDsConnection,
		error)
}
