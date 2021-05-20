package resolver

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageResolver struct{}

func NewMessageResolver() *MessageResolver {
	return &MessageResolver{}
}

func (r *MessageResolver) User(ctx context.Context, obj *entity.Message) (*entity.User, error) {
	return nil, nil
}
