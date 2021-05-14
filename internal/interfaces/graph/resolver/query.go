package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/message"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/commander"
)

type QueryResolver struct {
	messageCommander commander.MessageCommander
	userCommander    commander.UserCommander
}

func NewQueryResolver(
	messageCommander commander.MessageCommander,
	userCommander commander.UserCommander,
) *QueryResolver {
	return &QueryResolver{
		messageCommander: messageCommander,
		userCommander:    userCommander,
	}
}

func (r *QueryResolver) Messages(ctx context.Context) ([]*message.Message, error) {
	mm, err := r.messageCommander.Messages(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Message Commander] message: %w", err)
	}

	return mm, nil
}

func (r *QueryResolver) Users(ctx context.Context) ([]string, error) {
	users, err := r.userCommander.Users(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Commander] users: %w", err)
	}

	return users, nil
}
