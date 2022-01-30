package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
)

type UserResolver struct {
	userLoader loader.UserLoader
}

func NewUserResolver(userLoader loader.UserLoader) *UserResolver {
	return &UserResolver{
		userLoader: userLoader,
	}
}

func (r *UserResolver) Friends(
	ctx context.Context,
	obj *entity.User,
	first int,
	after entity.ID,
	sortBy entity.FriendsSortByType,
	sortOrder entity.SortOrderType,
) (*entity.FriendsConnection, error) {
	idsCon, err := r.userLoader.LoadFriendIDs(ctx, entity.FriendsQueryInput{
		UserID: obj.ID,
		ListQueryInput: entity.ListQueryInput{
			First:     first,
			After:     after,
			SortBy:    sortBy,
			SortOrder: sortOrder,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("load friend ids: %w", err)
	}

	if idsCon == nil {
		return nil, fmt.Errorf("load friend ids: %w",
			fmt.Errorf("ids connection nil"))
	}

	friendsEdges := make([]*entity.FriendsEdge, 0, len(idsCon.Edges))

	for _, edge := range idsCon.Edges {
		user, err := r.userLoader.LoadUser(ctx, edge.Node)
		if err != nil {
			return nil, fmt.Errorf("load user: %w", err)
		}

		friendsEdges = append(friendsEdges, &entity.FriendsEdge{
			Cursor: user.ID,
			Node:   user,
		})
	}

	return &entity.FriendsConnection{
		PageInfo:   idsCon.PageInfo,
		Edges:      friendsEdges,
		TotalCount: idsCon.TotalCount,
	}, nil
}

func (r *UserResolver) Conversations(
	ctx context.Context,
	obj *entity.User,
	first int,
	after entity.ID,
) (*entity.ConversationsConnection, error) {
	// TODO
	return nil, nil
}
