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
	transactor        repository.Transactor
}

func NewMessageUsecase(
	userRepository repository.UserRepository,
	messageRepository repository.MessageRepository,
	transactor repository.Transactor,
) *MessageUsecase {
	return &MessageUsecase{
		userRepository:    userRepository,
		messageRepository: messageRepository,
		transactor:        transactor,
	}
}

func (u *MessageUsecase) PostMessage(
	ctx context.Context,
	conversationID entity.ID,
	msgType entity.MessageType,
	senderID entity.ID,
	text string,
) (*entity.Message, error) {
	txCtx, err := u.transactor.Begin(ctx)
	if err != nil {
		return nil, errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("begin transaction: %w", err))
	}

	message, err := u.messageRepository.CreateMessageWithTransaction(txCtx, conversationID,
		msgType, senderID, text)
	if err != nil {
		return nil, errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("create message: %w", err))
	}

	if err := u.transactor.Commit(txCtx); err != nil {
		return nil, errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("commit transaction: %w", err))
	}

	// skip error when fanout message
	u.messageRepository.FanoutMessage(ctx, message)

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
	fail := func(err error) (*entity.Conversation, error) {
		return nil, fmt.Errorf("CreateNewConversation: %w", err)
	}

	txCtx, err := u.transactor.Begin(ctx)
	if err != nil {
		return fail(errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("begin transaction: %w", err)))
	}

	conversationID, err := u.messageRepository.CreateConversationWithTransaction(txCtx,
		creatorID, conversationTitle, conversationType, recipentIDs)
	if err != nil {
		return fail(errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("create conversation: %w", err)))
	}

	if text != nil {
		_, err := u.messageRepository.CreateMessageWithTransaction(txCtx, *conversationID,
			entity.MessageTypeText, creatorID, *text)
		if err != nil {
			return nil, errorHandlerWithTransaction(txCtx, u.transactor,
				fmt.Errorf("create message: %w", err))
		}
	}

	cc, err := u.messageRepository.FindConversationsByIDsWithTransaction(txCtx,
		[]entity.ID{*conversationID})
	if err != nil {
		return nil, errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("find conversation: %w", err))
	}

	if len(cc) == 0 {
		return nil, errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("conversation not exist"))
	}

	if err := u.transactor.Commit(txCtx); err != nil {
		return nil, errorHandlerWithTransaction(txCtx, u.transactor,
			fmt.Errorf("commit transaction: %w", err))
	}

	return cc[0], nil
}

func (u *MessageUsecase) MessagePosted(ctx context.Context) (
	<-chan *entity.Message, error) {
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

func (u *MessageUsecase) AllMessagesByConversationIDs(ctx context.Context,
	conversationIDs []entity.ID) (
	map[entity.ID][]*entity.Message, error) {
	res, err := u.messageRepository.FindAllMessagesInConversations(ctx,
		conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("find messages in conversations: %w", err)
	}

	return res, nil
}

func (u *MessageUsecase) ConversationByIDs(ctx context.Context,
	conversationIDs []entity.ID) ([]*entity.Conversation,
	error) {
	res, err := u.messageRepository.FindConversationsByIDs(ctx, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("find conversations: %w", err)
	}

	return res, nil
}

func (u *MessageUsecase) Conversations(ctx context.Context) (
	[]*entity.Conversation, error) {
	user, err := u.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	// TODO:
	return nil, nil
}

func (u *MessageUsecase) GetConversationIDsFromUserIDs(
	ctx context.Context,
	inputs []entity.RelayQueryInput,
) (map[entity.ID]*entity.IDsConnection, error) {
	users, err := u.messageRepository.FindConversationIDsFromUserIDs(ctx, inputs)
	if err != nil {
		return nil, fmt.Errorf("find conversation ids from user ids: %w", err)
	}

	return users, nil
}

func (u *MessageUsecase) GetParticipantsInConversations(ctx context.Context,
	conversationIDs []entity.ID) (map[entity.ID][]*entity.User, error) {
	participants, err := u.messageRepository.FindParticipantsInConversations(ctx,
		conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("find participants in conversations: %w", err)
	}

	return participants, nil
}

func (u *MessageUsecase) MessagesInConversations(ctx context.Context,
	inputs []entity.RelayQueryInput,
) (map[entity.ID]*entity.ConversationMessagesConnection, error) {
	res, err := u.messageRepository.FindMessagesInConversations(ctx, inputs)
	if err != nil {
		return nil, fmt.Errorf("find messages in conversations: %w", err)
	}

	return res, nil
}
