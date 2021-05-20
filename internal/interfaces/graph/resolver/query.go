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
		return nil, fmt.Errorf("[Message Commander] message: %w", err)
	}

	return mm, nil
}

func (r *QueryResolver) Users(ctx context.Context) ([]*entity.User, error) {
	users, err := r.userUsecase.Users(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Commander] users: %w", err)
	}

	return users, nil
}
