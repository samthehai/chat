package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/samthehai/chat/internal/domain/message"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
)

const messagesKey = "messages"

type MessageRepository struct {
	cacher   external.Cacher
	msgChans map[string]chan *message.Message
	mutex    sync.Mutex
}

func NewMessageRepository(
	cacher external.Cacher,
) *MessageRepository {
	return &MessageRepository{
		cacher:   cacher,
		msgChans: map[string]chan *message.Message{},
		mutex:    sync.Mutex{},
	}
}

func (r *MessageRepository) PostMessage(_ context.Context, msg *message.Message) error {
	mj, _ := json.Marshal(msg)
	if err := r.cacher.LPush(messagesKey, mj); err != nil {
		return fmt.Errorf("[cacher] lpush: %w", err)
	}

	r.mutex.Lock()
	for _, ch := range r.msgChans {
		ch <- msg
	}
	r.mutex.Unlock()

	return nil
}

func (r *MessageRepository) Messages(_ context.Context) ([]*message.Message, error) {
	values, err := r.cacher.LRange(messagesKey, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("[cacher] messages: %w", err)
	}

	messages := []*message.Message{}
	for _, mj := range values {
		m := &message.Message{}
		if err := json.Unmarshal([]byte(mj), &m); err != nil {
			return nil, fmt.Errorf("[Unmarshal] messages: %w", err)
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func (s *MessageRepository) MessagePosted(ctx context.Context, user string) (<-chan *message.Message, error) {
	messages := make(chan *message.Message, 1)
	s.mutex.Lock()
	s.msgChans[user] = messages
	s.mutex.Unlock()

	go func() {
		<-ctx.Done()

		s.mutex.Lock()
		delete(s.msgChans, user)
		s.mutex.Unlock()
	}()

	return messages, nil
}
