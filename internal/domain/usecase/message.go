package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/domain/usecase/repository"
	"github.com/segmentio/ksuid"
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

func (c *MessageUsecase) PostMessage(ctx context.Context, text string) (*entity.Message, error) {
	user, err := c.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("[User Repository] user is nil")
	}

	m := &entity.Message{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Text:      text,
		UserID:    user.ID,
	}

	if err := c.messageRepository.PostMessage(ctx, m); err != nil {
		return nil, fmt.Errorf("[Message Repository] post message: %w", err)
	}

	return m, nil
}

func (c *MessageUsecase) Messages(ctx context.Context) ([]*entity.Message, error) {
	messages, err := c.messageRepository.Messages(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Message Repository] messages: %w", err)
	}

	return messages, nil
}

func (c *MessageUsecase) MessagePosted(ctx context.Context) (<-chan *entity.Message, error) {
	user, err := c.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("[User Repository] user is nil")
	}

	messages, err := c.messageRepository.MessagePosted(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("[Message Repository] message posted: %w", err)
	}

	return messages, nil
}
