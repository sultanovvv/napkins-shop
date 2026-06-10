package postgres

import (
	"context"

	"github.com/uptrace/bun"
)

func setTransactionInCtx(ctx context.Context, tx *bun.Tx, name ctxKey) context.Context {
	return context.WithValue(ctx, name, tx)
}

func getTransactionFromCtx(ctx context.Context, name ctxKey) (*bun.Tx, bool) {
	val := ctx.Value(name)
	if tx, ok := val.(*bun.Tx); ok {
		return tx, true
	}
	return nil, false
}
