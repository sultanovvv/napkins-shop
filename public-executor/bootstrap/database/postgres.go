package database

import (
	"fmt"

	"napkins-shop/public-executor/bootstrap/config"

	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const driverName = "pgx"

func NewConnect(cfg *config.DBConfig, log *zap.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		// https://github.com/jmoiron/sqlx/issues/787
		loggerOptions := []sqldblogger.Option{
			sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		}

		dblog := sqldblogger.OpenDriver(dsn, db.Driver(), zapadapter.New(log), loggerOptions...)
		db = sqlx.NewDb(dblog, driverName)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if errPing := db.Ping(); errPing != nil {
		return nil, errPing
	}

	return db, nil
}
