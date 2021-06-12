package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
)

type ParticipantResolver struct {
	userLoader loader.UserLoader
}

func NewParticipantResolver(
	userLoader loader.UserLoader,
) *ParticipantResolver {
	return &ParticipantResolver{
		userLoader: userLoader,
	}
}

func (r *ParticipantResolver) User(ctx context.Context, obj *entity.Participant) (*entity.User, error) {
	user, err := r.userLoader.LoadUser(ctx, obj.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return user, nil
}

func (r *ParticipantResolver) Conversation(ctx context.Context, obj *entity.Participant) (*entity.Conversation, error) {
	return nil, nil
}
