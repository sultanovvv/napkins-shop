package app

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"napkins-shop/backoffice-executor/configs"
	httpentry "napkins-shop/backoffice-executor/internal/entrypoints/http"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/docs"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/health"
	"napkins-shop/backoffice-executor/internal/middlewares"
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
		fx.Provide(health.NewHealthHandler),
		fx.Invoke(middlewares.Register, httpentry.Register, docs.Register),
	)
}
