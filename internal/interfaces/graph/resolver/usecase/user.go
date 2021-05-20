package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type UserUsecase interface {
	Users(ctx context.Context) ([]*entity.User, error)
	UserJoined(ctx context.Context) (<-chan *entity.User, error)
	LoginWithFacebook(ctx context.Context) (*entity.User, error)
}
