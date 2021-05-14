package commander

import (
	"context"

	"github.com/samthehai/chat/internal/domain/message"
)

type MessageCommander interface {
	PostMessage(ctx context.Context, user string, text string) (*message.Message, error)
	Messages(ctx context.Context) ([]*message.Message, error)
	MessagePosted(ctx context.Context, user string) (<-chan *message.Message, error)
}
