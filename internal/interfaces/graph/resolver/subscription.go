package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

type SubscriptionResolver struct {
	messageUsecase usecase.MessageUsecase
	userUsecase    usecase.UserUsecase
}

func NewSubscriptionResolver(
	messageUsecase usecase.MessageUsecase,
	userUsecase usecase.UserUsecase,
) *SubscriptionResolver {
	return &SubscriptionResolver{
		messageUsecase: messageUsecase,
		userUsecase:    userUsecase,
	}
}

func (r *SubscriptionResolver) MessagePosted(ctx context.Context) (<-chan *entity.Message, error) {
	messages, err := r.messageUsecase.MessagePosted(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Message commander] message posted: %w", err)
	}

	return messages, nil
}

func (r *SubscriptionResolver) UserJoined(ctx context.Context) (<-chan *entity.User, error) {
	users, err := r.userUsecase.UserJoined(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User commander] user joined: %w", err)
	}

	return users, nil
}
