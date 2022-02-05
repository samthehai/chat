package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageUsecase interface {
	AllMessagesByConversationIDs(ctx context.Context, conversationIDs []entity.ID) (
		map[entity.ID][]*entity.Message, error)
	ConversationByIDs(ctx context.Context, conversationIDs []entity.ID) (
		[]*entity.Conversation, error)
	GetConversationIDsFromUserIDs(ctx context.Context,
		inputs []entity.RelayQueryInput) (map[entity.ID]*entity.IDsConnection, error)
	GetParticipantsInConversations(ctx context.Context,
		conversationIDs []entity.ID) (map[entity.ID][]*entity.User, error)
	MessagesInConversations(ctx context.Context,
		inputs []entity.RelayQueryInput,
	) (map[entity.ID]*entity.ConversationMessagesConnection, error)
}
