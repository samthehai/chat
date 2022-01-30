package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
)

type MessageResolver struct {
	userLoader         loader.UserLoader
	conversationLoader loader.ConversationLoader
}

func NewMessageResolver(
	userLoader loader.UserLoader,
	conversationLoader loader.ConversationLoader,
) *MessageResolver {
	return &MessageResolver{
		userLoader:         userLoader,
		conversationLoader: conversationLoader,
	}
}

func (r *MessageResolver) Sender(
	ctx context.Context,
	obj *entity.Message,
) (*entity.User, error) {
	sender, err := r.userLoader.LoadUser(ctx, obj.SenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return sender, nil
}

func (r *MessageResolver) Conversation(
	ctx context.Context,
	obj *entity.Message,
) (*entity.Conversation, error) {
	c, err := r.conversationLoader.LoadConversation(ctx, obj.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}

	return c, nil
}
