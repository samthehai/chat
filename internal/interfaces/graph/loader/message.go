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
			messageFetcher.MessagesByConversationIDs,
		),
	}
}

func (l *MessageLoader) LoadMessagesInConversation(
	ctx context.Context,
	conversationID entity.ID,
) ([]*entity.Message, error) {
	raw, err := l.messagesInConversationLoader.Load(ctx, conversationID)()
	if err != nil {
		return nil, fmt.Errorf("load messages in conversation: id=%v, %w", conversationID, err)
	}

	if raw == nil {
		return nil, nil
	}

	return raw.([]*entity.Message), nil
}

func newMessagesInConversationLoader(
	fetchFunc func(ctx context.Context, conversationIDs []entity.ID) (map[entity.ID][]*entity.Message, error),
) *dataloader.Loader {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			ids := getIDsFromKeys(keys)

			mapMsgs, err := fetchFunc(ctx, ids)
			if err != nil {
				return fillUpResultsWithError(len(keys), err)
			}

			results := make([]*dataloader.Result, 0, len(keys))
			for _, id := range ids {
				d := mapMsgs[id]

				results = append(results, &dataloader.Result{
					Data:  d,
					Error: err,
				})
			}

			return results
		},
	)
}
