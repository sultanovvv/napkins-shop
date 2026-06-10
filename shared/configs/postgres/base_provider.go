package postgres

import (
	"context"

	"github.com/uptrace/bun"
)

// IBaseProvider — точка доступа к БД для репозитория.
// Conn(ctx) возвращает текущую транзакцию, если она положена в контекст
// одноимённым TransactionProvider'ом; иначе — корневой *bun.DB.
// ConnWithoutTX игнорирует контекст (для запросов, которые осознанно
// не должны попадать в транзакцию — health-check, чтение из кэшируемой таблицы).
type IBaseProvider interface {
	Conn(ctx context.Context) bun.IDB
	ConnWithoutTX() bun.IDB
}

type ctxKey string

type baseProvider struct {
	db   *bun.DB
	name ctxKey
}

func (p baseProvider) Conn(ctx context.Context) bun.IDB {
	tx, ok := getTransactionFromCtx(ctx, p.name)
	if !ok {
		return p.db
	}
	return tx
}

func (p baseProvider) ConnWithoutTX() bun.IDB {
	return p.db
}

// NewNamedBaseProvider — имя нужно для разводки нескольких независимых
// connection-пулов (например, core-db и read-replica) — каждый кладёт
// свою транзакцию в контекст под своим ключом.
func NewNamedBaseProvider(db *bun.DB, name string) IBaseProvider {
	return baseProvider{
		db:   db,
		name: ctxKey(name),
	}
}
