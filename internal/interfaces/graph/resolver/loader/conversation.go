package loader

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type ConversationLoader interface {
	LoadConversation(
		ctx context.Context,
		conversationID entity.ID,
	) (*entity.Conversation, error)
}
