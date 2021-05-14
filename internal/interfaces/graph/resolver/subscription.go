package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/message"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/commander"
)

type SubscriptionResolver struct {
	messageCommander commander.MessageCommander
	userCommander    commander.UserCommander
}

func NewSubscriptionResolver(
	messageCommander commander.MessageCommander,
	userCommander commander.UserCommander,
) *SubscriptionResolver {
	return &SubscriptionResolver{
		messageCommander: messageCommander,
		userCommander:    userCommander,
	}
}

func (r *SubscriptionResolver) MessagePosted(ctx context.Context, user string) (<-chan *message.Message, error) {
	messages, err := r.messageCommander.MessagePosted(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("[Message commander] message posted: %w", err)
	}

	return messages, nil
}

func (r *SubscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	users, err := r.userCommander.UserJoined(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("[User commander] user joined: %w", err)
	}

	return users, nil
}
