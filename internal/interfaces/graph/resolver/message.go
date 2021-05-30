package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

type MessageResolver struct {
	userUsecase usecase.UserUsecase
}

func NewMessageResolver(
	userUsecase usecase.UserUsecase,
) *MessageResolver {
	return &MessageResolver{
		userUsecase: userUsecase,
	}
}

func (r *MessageResolver) User(ctx context.Context, obj *entity.Message) (*entity.User, error) {
	user, err := r.userUsecase.User(ctx, obj.UserID)
	if err != nil {
		return nil, fmt.Errorf("[User Usecase] user: %w", err)
	}

	return user, nil
}
