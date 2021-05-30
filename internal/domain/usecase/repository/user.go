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
	FindFriends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error)
	FindUser(ctx context.Context, id entity.ID) (*entity.User, error)
}
