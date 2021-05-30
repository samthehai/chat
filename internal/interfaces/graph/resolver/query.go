package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

type QueryResolver struct {
	messageUsecase usecase.MessageUsecase
	userUsecase    usecase.UserUsecase
}

func NewQueryResolver(
	messageUsecase usecase.MessageUsecase,
	userUsecase usecase.UserUsecase,
) *QueryResolver {
	return &QueryResolver{
		messageUsecase: messageUsecase,
		userUsecase:    userUsecase,
	}
}

func (r *QueryResolver) Messages(ctx context.Context) ([]*entity.Message, error) {
	mm, err := r.messageUsecase.Messages(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Message Usecase] message: %w", err)
	}

	return mm, nil
}

func (r *QueryResolver) Friends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error) {
	friends, err := r.userUsecase.Friends(ctx, first, after, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("[User Usecase] friends: %w", err)
	}

	return friends, nil
}
