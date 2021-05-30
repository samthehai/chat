package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type UserUsecase interface {
	Friends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error)
	UserJoined(ctx context.Context) (<-chan *entity.User, error)
	CurrentUser(ctx context.Context) (*entity.User, error)
	User(ctx context.Context, id entity.ID) (*entity.User, error)
}
