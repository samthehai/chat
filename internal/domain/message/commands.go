package message

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

type MessageCommander struct {
	userRepository    UserRepository
	messageRepository MessageRepository
}

func NewMessageCommander(
	userRepository UserRepository,
	messageRepository MessageRepository,
) *MessageCommander {
	return &MessageCommander{
		userRepository:    userRepository,
		messageRepository: messageRepository,
	}
}

func (c *MessageCommander) PostMessage(ctx context.Context, user string, text string) (*Message, error) {
	if err := c.userRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("[User Repository] create user: %w", err)
	}

	m := &Message{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Text:      text,
		User:      user,
	}

	if err := c.messageRepository.PostMessage(ctx, m); err != nil {
		return nil, fmt.Errorf("[Message Repository] post message: %w", err)
	}

	return m, nil
}

func (c *MessageCommander) Messages(ctx context.Context) ([]*Message, error) {
	messages, err := c.messageRepository.Messages(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Message Repository] messages: %w", err)
	}

	return messages, nil
}

func (c *MessageCommander) MessagePosted(ctx context.Context, user string) (<-chan *Message, error) {
	if err := c.userRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("[User Repository] create user: %w", err)
	}

	messages, err := c.messageRepository.MessagePosted(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("[Message Repository] message posted: %w", err)
	}

	return messages, nil
}
