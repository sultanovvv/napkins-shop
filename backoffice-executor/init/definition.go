package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"napkins-shop/backoffice-executor/configs"
	httpentry "napkins-shop/backoffice-executor/internal/entrypoints/http"
	categoryentry "napkins-shop/backoffice-executor/internal/entrypoints/http/category"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/docs"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/errorhandler"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/health"
	"napkins-shop/backoffice-executor/internal/middlewares"
	categoryuc "napkins-shop/backoffice-executor/internal/usecases/category"
	core_db "shared/providers/core-db"
)

func Infrastructure(a *App) fx.Option {
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
		fx.Provide(func() *configs.AppConfig { return a.appConf }),
		fx.Provide(func() *zap.Logger { return a.logger }),
		fx.Provide(func() *echo.Echo { return a.server }),
	)
}

func Providers() fx.Option {
	return fx.Options(
		fx.Provide(
			provideDatabase,
			core_db.NewBaseProvider,
			core_db.NewTransactionProvider,
			core_db.NewBackofficeCategoryProvider,
			categoryuc.NewUseCase,
			categoryentry.NewCategoryHandler,
			health.NewHealthHandler,
			errorhandler.New,
		),
		fx.Invoke(errorhandler.Register, middlewares.Register, httpentry.Register, docs.Register),
	)
}

// Один Connection используется провайдером и менеджером транзакций.
func provideDatabase(lifecycle fx.Lifecycle, cfg *configs.AppConfig) (*core_db.Connection, error) {
	db, err := core_db.InitDatabase(*cfg.Postgres)
	if err != nil {
		return nil, err
	}
	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error { return db.Close() },
	})
	return core_db.NewConnection(db), nil
}
