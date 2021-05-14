package user

import (
	"context"
	"fmt"
)

type UserCommander struct {
	userRepository UserRepository
}

func NewUserCommander(
	userRepository UserRepository,
) *UserCommander {
	return &UserCommander{
		userRepository: userRepository,
	}
}

func (c *UserCommander) CreateUser(ctx context.Context, user string) error {
	if err := c.userRepository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("[User Repository] create user: %w", err)
	}

	return nil
}

func (c *UserCommander) Users(ctx context.Context) ([]string, error) {
	users, err := c.userRepository.Users(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] users: %w", err)
	}

	return users, nil
}

func (c *UserCommander) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	if err := c.userRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("[User Repository] create user: %w", err)
	}

	users, err := c.userRepository.UserJoined(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] user joined: %w", err)
	}

	return users, nil
}
