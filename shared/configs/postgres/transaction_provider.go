package postgres

import (
	"context"
	"log/slog"

	"github.com/uptrace/bun"
)

type TransactionCallback func(ctx context.Context) error

// ITransactionProvider оборачивает Run в bun.BeginTx. Все BaseProvider'ы
// с тем же name внутри cb получают эту транзакцию через контекст.
type ITransactionProvider interface {
	Run(ctx context.Context, cb TransactionCallback) error
}

type transactionProvider struct {
	db   *bun.DB
	name ctxKey
}

func (p *transactionProvider) Run(ctx context.Context, cb TransactionCallback) error {
	tx, ok := getTransactionFromCtx(ctx, p.name)

	if ok {
		// Уже внутри транзакции — открываем nested через SAVEPOINT.
		ntx, err := tx.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		tx = &ntx
	} else {
		ntx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		tx = &ntx
	}

	ctx = setTransactionInCtx(ctx, tx, p.name)

	isCommitted := false
	defer func() {
		if !isCommitted {
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "rollback failed", "error", err)
			}
		}
	}()

	if err := cb(ctx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	isCommitted = true
	return nil
}

func NewNamedTransactionProvider(db *bun.DB, name string) ITransactionProvider {
	return &transactionProvider{
		db:   db,
		name: ctxKey(name),
	}
}
