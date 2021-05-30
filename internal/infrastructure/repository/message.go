package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
)

const messagesKey = "messages"

type MessageRepository struct {
	cacher   external.Cacher
	msgChans map[entity.ID]chan *entity.Message
	mutex    sync.Mutex
}

func NewMessageRepository(
	cacher external.Cacher,
) *MessageRepository {
	return &MessageRepository{
		cacher:   cacher,
		msgChans: map[entity.ID]chan *entity.Message{},
		mutex:    sync.Mutex{},
	}
}

func (r *MessageRepository) PostMessage(_ context.Context, msg *entity.Message) error {
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

func (r *MessageRepository) Messages(_ context.Context) ([]*entity.Message, error) {
	values, err := r.cacher.LRange(messagesKey, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("[cacher] messages: %w", err)
	}

	messages := []*entity.Message{}
	for _, mj := range values {
		m := &entity.Message{}
		if err := json.Unmarshal([]byte(mj), &m); err != nil {
			return nil, fmt.Errorf("[Unmarshal] messages: %w", err)
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func (s *MessageRepository) MessagePosted(ctx context.Context, input entity.User) (<-chan *entity.Message, error) {
	messages := make(chan *entity.Message, 1)
	s.mutex.Lock()
	s.msgChans[input.ID] = messages
	s.mutex.Unlock()

	go func() {
		<-ctx.Done()

		s.mutex.Lock()
		delete(s.msgChans, input.ID)
		s.mutex.Unlock()
	}()

	return messages, nil
}
