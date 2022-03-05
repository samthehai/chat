package usecase

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/usecase/repository"
)

func errorHandlerWithTransaction(ctx context.Context,
	transactor repository.Transactor, err error) error {
	if rbErr := transactor.Rollback(ctx); rbErr != nil {
		return fmt.Errorf("rollback err=%v: %w", err, rbErr)
	}

	return err
}
