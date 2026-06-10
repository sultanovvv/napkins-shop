package core_db

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"shared/configs/postgres"
)

// Connection — именованный wrapper над *bun.DB. Имя нужно, чтобы
// BaseProvider и TransactionProvider знали, под каким ключом класть
// транзакцию в контекст. Все сервисы, которые ходят в одну и ту же
// базу, используют name="core-db".
type Connection struct {
	*bun.DB
	name string
}

func NewConnection(db *bun.DB) *Connection {
	return &Connection{
		DB:   db,
		name: "core-db",
	}
}

func NewBaseProvider(connection *Connection) postgres.IBaseProvider {
	return postgres.NewNamedBaseProvider(connection.DB, connection.name)
}

type ITransactionProvider interface {
	postgres.ITransactionProvider
}

func NewTransactionProvider(connection *Connection) ITransactionProvider {
	return postgres.NewNamedTransactionProvider(connection.DB, connection.name)
}

// InitDatabase создаёт *bun.DB поверх pgx-stdlib. pgx используется как
// драйвер (database/sql-совместимый), bun — как ORM. Соединение пингуется
// перед возвратом, чтобы упасть пораньше, если конфиг кривой.
func InitDatabase(cfg postgres.Config) (*bun.DB, error) {
	connString := "postgres://" + cfg.User + ":" + cfg.Password +
		"@" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + "/" + cfg.Name + "?sslmode=disable"

	connCfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	sqlDB := stdlib.OpenDB(*connCfg)

	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.WithEnabled(cfg.Debug),
	))

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database %s@%s:%d/%s: %w",
			cfg.User, cfg.Host, cfg.Port, cfg.Name, err)
	}

	return db, nil
}
