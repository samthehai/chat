package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/message"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/commander"
)

type MutationResolver struct {
	messageCommander commander.MessageCommander
}

func NewMutationResolver(
	messageCommander commander.MessageCommander,
) *MutationResolver {
	return &MutationResolver{
		messageCommander: messageCommander,
	}
}

func (r *MutationResolver) PostMessage(ctx context.Context, user string, text string) (*message.Message, error) {
	m, err := r.messageCommander.PostMessage(ctx, user, text)
	if err != nil {
		return nil, fmt.Errorf("[Message Commander] post message: %w", err)
	}

	return m, nil
}
