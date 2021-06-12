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
}
