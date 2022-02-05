package loader

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type UserLoader interface {
	LoadUser(
		ctx context.Context,
		userID entity.ID,
	) (*entity.User, error)
	LoadFriendIDs(ctx context.Context,
		input entity.RelayQueryInput) (*entity.IDsConnection, error)
}
