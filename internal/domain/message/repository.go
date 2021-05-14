package message

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user string) error
}

type MessageRepository interface {
	PostMessage(ctx context.Context, msg *Message) error
	Messages(ctx context.Context) ([]*Message, error)
	MessagePosted(ctx context.Context, user string) (<-chan *Message, error)
}
