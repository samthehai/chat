package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageUsecase interface {
	PostMessage(ctx context.Context, text string) (*entity.Message, error)
	Messages(ctx context.Context) ([]*entity.Message, error)
	MessagePosted(ctx context.Context) (<-chan *entity.Message, error)
}
