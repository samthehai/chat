package loader

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageLoader interface {
	LoadMessagesInConversation(
		ctx context.Context,
		conversationID entity.ID,
	) ([]*entity.Message, error)
}
