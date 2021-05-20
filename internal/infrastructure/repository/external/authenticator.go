package external

import (
	"context"

	"github.com/samthehai/chat/internal/domain/entity"
)

type Authenticator interface {
	GetAuthTokenFromContext(ctx context.Context) (*entity.AuthToken, error)
}
