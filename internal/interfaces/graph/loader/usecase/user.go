package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type UserUsecase interface {
	Users(ctx context.Context, ids []entity.ID) ([]*entity.User, error)
}
