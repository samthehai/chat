package resolver

import (
	"context"

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

func (r *QueryResolver) Me(ctx context.Context) (*entity.User, error) {
	return r.userUsecase.Me(ctx)
}
