package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageUsecase interface {
	MessagesByConversationIDs(ctx context.Context, conversationIDs []entity.ID) (
		map[entity.ID][]*entity.Message, error)
	ConversationByIDs(ctx context.Context, conversationIDs []entity.ID) (
		[]*entity.Conversation, error)
	GetConversationIDsFromUserIDs(ctx context.Context,
		inputs []entity.UserQueryInput) (map[entity.ID]*entity.IDsConnection, error)
}
