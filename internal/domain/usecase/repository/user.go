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
	Users(ctx context.Context) ([]*entity.User, error)
	UserJoined(ctx context.Context, user entity.User) (<-chan *entity.User, error)
}
