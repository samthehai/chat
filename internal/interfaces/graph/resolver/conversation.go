package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
)

type ConversationResolver struct {
	messageLoader loader.MessageLoader
	userLoader    loader.UserLoader
}

func NewConversationResolver(
	messageLoader loader.MessageLoader,
	userLoader loader.UserLoader,
) *ConversationResolver {
	return &ConversationResolver{
		messageLoader: messageLoader,
		userLoader:    userLoader,
	}
}

func (r *ConversationResolver) Creator(ctx context.Context, obj *entity.Conversation) (*entity.User, error) {
	creator, err := r.userLoader.LoadUser(ctx, *obj.CreatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return creator, nil
}

func (r *ConversationResolver) Messages(ctx context.Context, obj *entity.Conversation) ([]*entity.Message, error) {
	messages, err := r.messageLoader.LoadMessagesInConversation(ctx, obj.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load messages in conversation: %w", err)
	}

	return messages, nil
}