package repository

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageRepository interface {
	PostMessage(ctx context.Context, msg *entity.Message) error
	Messages(ctx context.Context) ([]*entity.Message, error)
	MessagePosted(ctx context.Context, user entity.User) (<-chan *entity.Message, error)
}
