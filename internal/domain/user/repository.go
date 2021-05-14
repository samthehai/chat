package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user string) error
	Users(ctx context.Context) ([]string, error)
	UserJoined(ctx context.Context, user string) (<-chan string, error)
}
