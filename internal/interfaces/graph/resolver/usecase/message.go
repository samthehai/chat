package usecase

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageUsecase interface {
	PostMessage(
		ctx context.Context,
		conversationID entity.ID,
		msgType entity.MessageType,
		senderID entity.ID,
		text string,
	) (*entity.Message, error)
	CreateNewConversation(
		ctx context.Context,
		creatorID entity.ID,
		conversationTitle string,
		conversationType entity.ConversationType,
		recipentIDs []entity.ID,
		text *string,
	) (*entity.Conversation, error)
	MessagePosted(ctx context.Context) (<-chan *entity.Message, error)
}
