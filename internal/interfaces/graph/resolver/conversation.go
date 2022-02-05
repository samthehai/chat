package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
)

type ConversationResolver struct {
	messageLoader      loader.MessageLoader
	userLoader         loader.UserLoader
	conversationLoader loader.ConversationLoader
}

func NewConversationResolver(
	messageLoader loader.MessageLoader,
	userLoader loader.UserLoader,
	conversationLoader loader.ConversationLoader,
) *ConversationResolver {
	return &ConversationResolver{
		messageLoader:      messageLoader,
		userLoader:         userLoader,
		conversationLoader: conversationLoader,
	}
}

func (r *ConversationResolver) Creator(
	ctx context.Context,
	obj *entity.Conversation,
) (*entity.User, error) {
	creator, err := r.userLoader.LoadUser(ctx, *obj.CreatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return creator, nil
}

func (r *ConversationResolver) Messages(
	ctx context.Context,
	obj *entity.Conversation,
	first int,
	after entity.ID,
	sortBy entity.MessagesSortByType,
	sortOrder entity.SortOrderType,
) (*entity.ConversationMessagesConnection, error) {
	msgs, err := r.messageLoader.LoadMessagesInConversation(ctx,
		entity.RelayQueryInput{
			KeyID: obj.ID,
			ListQueryInput: entity.ListQueryInput{
				First:     first,
				After:     after,
				SortBy:    string(sortBy),
				SortOrder: sortOrder,
			},
		})
	if err != nil {
		return nil, fmt.Errorf("load messages in conversation: %w", err)
	}

	return msgs, nil
}

func (r *ConversationResolver) Participants(
	ctx context.Context,
	obj *entity.Conversation,
) ([]*entity.User, error) {
	pp, err := r.conversationLoader.LoadParticipantsInConversation(ctx, obj.ID)
	if err != nil {
		return nil, fmt.Errorf("load participants in conversation: %w", err)
	}

	return pp, nil
}
