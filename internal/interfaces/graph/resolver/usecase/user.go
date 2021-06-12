package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type UserUsecase interface {
	GetUserFromContext(ctx context.Context) (*entity.User, error)
	Friends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error)
	UserJoined(ctx context.Context) (<-chan *entity.User, error)
	Login(ctx context.Context) (*entity.User, error)
	Me(ctx context.Context) (*entity.User, error)
}
