package loader

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageLoader interface {
	LoadMessagesInConversation(
		ctx context.Context,
		input entity.RelayQueryInput,
	) (*entity.ConversationMessagesConnection, error)
}
