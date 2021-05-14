package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/samthehai/chat/internal/infrastructure/repository/external"
)

const usersKey = "users"

type UserRepository struct {
	cacher    external.Cacher
	userChans map[string]chan string
	mutex     sync.Mutex
}

func NewUserRepository(
	cacher external.Cacher,
) *UserRepository {
	return &UserRepository{
		cacher:    cacher,
		userChans: map[string]chan string{},
		mutex:     sync.Mutex{},
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user string) error {
	if err := r.cacher.SAdd(usersKey, []byte(user)); err != nil {
		return err
	}

	r.mutex.Lock()
	for _, ch := range r.userChans {
		ch <- user
	}
	r.mutex.Unlock()

	return nil
}

func (r *UserRepository) Users(ctx context.Context) ([]string, error) {
	users, err := r.cacher.SMembers(usersKey)
	if err != nil {
		return nil, fmt.Errorf("cacher] smembers: %w", err)
	}

	return users, nil
}

func (r *UserRepository) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	users := make(chan string, 1)

	r.mutex.Lock()
	r.userChans[user] = users
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()

		r.mutex.Lock()
		delete(r.userChans, user)
		r.mutex.Unlock()
	}()

	return users, nil
}
