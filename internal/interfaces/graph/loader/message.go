package loader

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/loader/usecase"
)

type MessageLoader struct {
	messagesInConversationLoader *dataloader.Loader
}

func NewMessageLoader(
	messageFetcher usecase.MessageUsecase,
) *MessageLoader {
	return &MessageLoader{
		messagesInConversationLoader: newMessagesInConversationLoader(
			messageFetcher.MessagesInConversations,
		),
	}
}

func (l *MessageLoader) LoadMessagesInConversation(
	ctx context.Context,
	input entity.RelayQueryInput,
) (*entity.ConversationMessagesConnection, error) {
	raw, err := l.messagesInConversationLoader.Load(ctx, input)()
	if err != nil {
		return nil, fmt.Errorf("load messages in conversation: input=%v, %w", input, err)
	}

	return raw.(*entity.ConversationMessagesConnection), nil
}

func newMessagesInConversationLoader(
	fetchFunc func(ctx context.Context, inputs []entity.RelayQueryInput) (
		map[entity.ID]*entity.ConversationMessagesConnection, error),
) *dataloader.Loader {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			inputs := make([]entity.RelayQueryInput, 0, len(keys))

			for _, key := range keys {
				inputs = append(inputs, key.Raw().(entity.RelayQueryInput))
			}

			mapMsgs, err := fetchFunc(ctx, inputs)
			if err != nil {
				return fillUpResultsWithError(len(keys), err)
			}

			results := make([]*dataloader.Result, 0, len(keys))
			for _, input := range inputs {
				d := mapMsgs[input.KeyID]

				results = append(results, &dataloader.Result{
					Data:  d,
					Error: err,
				})
			}

			return results
		},
	)
}
