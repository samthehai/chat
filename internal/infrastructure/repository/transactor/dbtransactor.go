package transactor

import (
	"context"
	"database/sql"
	"fmt"
)

type transactionKeyType string

const transactionKey transactionKeyType = "TRANSACTION_KEY"

type DBTransactor struct {
	db *sql.DB
}

func NewDBTransactor(db *sql.DB) *DBTransactor {
	return &DBTransactor{db: db}
}

func (rtx *DBTransactor) Begin(ctx context.Context) (context.Context, error) {
	tx, err := rtx.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return context.WithValue(ctx, transactionKey, tx), nil
}

func (rtx *DBTransactor) Commit(ctx context.Context) error {
	tx, ok := rtx.GetTransactionFromCtx(ctx)
	if !ok || tx == nil {
		return fmt.Errorf("can not get transation from ctx")
	}

	return tx.Commit()
}

func (rtx *DBTransactor) Rollback(ctx context.Context) error {
	tx, ok := rtx.GetTransactionFromCtx(ctx)
	if !ok || tx == nil {
		return fmt.Errorf("can not get transation from ctx")
	}

	return tx.Rollback()
}

func (rtx *DBTransactor) GetTransactionFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(transactionKey).(*sql.Tx)
	return tx, ok
}
