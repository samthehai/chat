package external

import (
	"context"
	"database/sql"
)

type Transactor interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetTransactionFromCtx(ctx context.Context) (*sql.Tx, bool)
}
