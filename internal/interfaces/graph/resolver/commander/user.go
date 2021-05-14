package commander

import "context"

type UserCommander interface {
	Users(ctx context.Context) ([]string, error)
	UserJoined(ctx context.Context, user string) (<-chan string, error)
}
