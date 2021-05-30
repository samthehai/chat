package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

type MutationResolver struct {
	messageUsecase usecase.MessageUsecase
	userUsecase    usecase.UserUsecase
}

func NewMutationResolver(
	messageUsecase usecase.MessageUsecase,
	userUsecase usecase.UserUsecase,
) *MutationResolver {
	return &MutationResolver{
		messageUsecase: messageUsecase,
		userUsecase:    userUsecase,
	}
}

func (r *MutationResolver) PostMessage(ctx context.Context, text string) (*entity.Message, error) {
	m, err := r.messageUsecase.PostMessage(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("[Message Commander] post message: %w", err)
	}

	return m, nil
}

func (r *MutationResolver) CurrentUser(ctx context.Context) (*entity.User, error) {
	return r.userUsecase.CurrentUser(ctx)
}
