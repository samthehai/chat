package usecase

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/domain/usecase/repository"
)

type MessageUsecase struct {
	userRepository    repository.UserRepository
	messageRepository repository.MessageRepository
}

func NewMessageUsecase(
	userRepository repository.UserRepository,
	messageRepository repository.MessageRepository,
) *MessageUsecase {
	return &MessageUsecase{
		userRepository:    userRepository,
		messageRepository: messageRepository,
	}
}

func (u *MessageUsecase) PostMessage(
	ctx context.Context,
	conversationID entity.ID,
	msgType entity.MessageType,
	senderID entity.ID,
	text string,
) (*entity.Message, error) {
	message, err := u.messageRepository.CreateMessage(ctx, conversationID, msgType, senderID, text)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return message, nil
}

func (u *MessageUsecase) CreateNewConversation(
	ctx context.Context,
	creatorID entity.ID,
	conversationTitle string,
	conversationType entity.ConversationType,
	recipentIDs []entity.ID,
	text *string,
) (*entity.Conversation, error) {
	conversationID, err := u.messageRepository.CreateConversation(ctx, creatorID, conversationTitle, conversationType, recipentIDs)
	if err != nil {
		return nil, fmt.Errorf("post message: %w", err)
	}

	if text != nil {
		_, err := u.messageRepository.CreateMessage(ctx, *conversationID, entity.MessageTypeText, creatorID, *text)
		if err != nil {
			return nil, fmt.Errorf("create message: %w", err)
		}
	}

	cc, err := u.messageRepository.FindConversations(ctx, []entity.ID{*conversationID})
	if err != nil {
		return nil, fmt.Errorf("find conversation: %w", err)
	}

	if len(cc) == 0 {
		return nil, fmt.Errorf("conversation not exist")
	}

	return cc[0], nil
}

func (u *MessageUsecase) MessagePosted(ctx context.Context) (<-chan *entity.Message, error) {
	user, err := u.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	messages, err := u.messageRepository.MessagePosted(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("message posted: %w", err)
	}

	return messages, nil
}

func (u *MessageUsecase) MessagesInConversation(ctx context.Context, conversationIDs []entity.ID) (map[entity.ID][]*entity.Message, error) {
	res, err := u.messageRepository.FindMessagesInConversations(ctx, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("find messages in conversations: %w", err)
	}

	return res, nil
}

func (u *MessageUsecase) Conversations(ctx context.Context, conversationIDs []entity.ID) ([]*entity.Conversation, error) {
	res, err := u.messageRepository.FindConversations(ctx, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("find conversations: %w", err)
	}

	return res, nil
}
