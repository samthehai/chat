package resolver

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/model"
)

type UserResolver struct {
}

func NewUserResolver() *UserResolver {
	return &UserResolver{}
}

func (r *UserResolver) Friends(
	ctx context.Context,
	obj *entity.User,
	first int,
	after entity.ID,
	sortBy entity.FriendsSortByType,
	sortOrder entity.SortOrderType,
) (*model.FriendsConnection, error) {
	return model.NewFriendsConnection()
}

func (r *UserResolver) Conversations(
	ctx context.Context,
	obj *entity.User,
	first int,
	after entity.ID,
) (*model.ConversationsConnection, error) {
	// TODO
	return nil, nil
}
