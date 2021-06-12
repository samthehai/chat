package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageUsecase interface {
	MessagesInConversation(ctx context.Context, conversationIDs []entity.ID) (map[entity.ID][]*entity.Message, error)
	Conversations(ctx context.Context, conversationIDs []entity.ID) ([]*entity.Conversation, error)
}
