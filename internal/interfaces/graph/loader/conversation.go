package loader

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/loader/usecase"
)

type ConversationLoader struct {
	conversationLoader *dataloader.Loader
}

func NewConversationLoader(
	messageUsecase usecase.MessageUsecase,
) *ConversationLoader {
	return &ConversationLoader{
		conversationLoader: newConversationLoader(
			messageUsecase.ConversationByIDs,
		),
	}
}

func (l *ConversationLoader) LoadConversation(
	ctx context.Context,
	conversationID entity.ID,
) (*entity.Conversation, error) {
	raw, err := l.conversationLoader.Load(ctx, conversationID)()
	if err != nil {
		return nil, fmt.Errorf("load conversation: id=%v, %w", conversationID, err)
	}

	return raw.(*entity.Conversation), nil
}

func newConversationLoader(
	fetchFunc func(ctx context.Context, conversationIDs []entity.ID) ([]*entity.Conversation, error),
) *dataloader.Loader {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			ids := getIDsFromKeys(keys)

			conversations, err := fetchFunc(ctx, ids)
			if err != nil {
				return fillUpResultsWithError(len(keys), err)
			}

			m := make(map[entity.ID]*entity.Conversation)
			for _, c := range conversations {
				m[c.ID] = c
			}

			results := make([]*dataloader.Result, 0, len(keys))
			for _, id := range ids {
				d := m[id]

				results = append(results, &dataloader.Result{
					Data:  d,
					Error: err,
				})
			}

			return results
		},
	)
}
