package repository

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type MessageRepository interface {
	CreateConversation(
		ctx context.Context,
		creatorID entity.ID,
		conversationTitle string,
		conversationType entity.ConversationType,
		recipentIDs []entity.ID,
	) (*entity.ID, error)
	FindConversationsByIDs(
		ctx context.Context,
		conversationIDs []entity.ID,
	) ([]*entity.Conversation, error)
	CreateMessage(
		ctx context.Context,
		conversationID entity.ID,
		msgType entity.MessageType,
		senderID entity.ID,
		msg string,
	) (*entity.Message, error)
	FindMessagesInConversations(
		ctx context.Context,
		conversationIDs []entity.ID,
	) (map[entity.ID][]*entity.Message, error)
	MessagePosted(
		ctx context.Context,
		user entity.User,
	) (<-chan *entity.Message, error)
	FanoutMessage(
		ctx context.Context,
		message *entity.Message,
	)
	FindConversationIDsFromUserIDs(ctx context.Context,
		inputs []entity.UserQueryInput) (map[entity.ID]*entity.IDsConnection, error)
}
