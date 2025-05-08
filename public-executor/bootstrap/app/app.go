package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brpaz/echozap"

	"napkins-shop/public-executor/bootstrap/config"
	"napkins-shop/public-executor/internal/definition"
	"napkins-shop/public-executor/internal/utils/logger"

	"github.com/labstack/echo/v4/middleware"
)

const shutdownTimeout = 10 * time.Second

type App struct {
	Logger *zap.Logger
	Config *config.AppConfig
	//DB     *sqlx.DB
}

func NewApp(configPath string) (*App, error) {
	config.InitViperByEnv(configPath)
	cfg := config.NewAppConfig()
	log := logger.NewLogger()

	var (
	//err error
	//db  *sqlx.DB
	)
	//if db, err = database.NewConnect(cfg.Database(), log); err != nil {
	//	log.Panic("Error connecting to Postgres", zap.String("host", cfg.Database().Host), zap.Error(err))
	//	return nil, err
	//}

	return &App{
		Logger: log,
		Config: cfg,
		//DB:     db,
	}, nil
}

func (app *App) Start() error {
	fxApp := fx.New(
		fx.Provide(func() *config.AppConfig { return app.Config }),
		fx.Provide(func() *zap.Logger { return app.Logger }),
		//fx.Provide(func() *sqlx.DB { return app.DB }),
		fx.Provide(func() *echo.Echo {
			e := echo.New()
			e.Use(echozap.ZapLogger(app.Logger))
			e.Use(middleware.RequestID())

			e.Use(middleware.Recover())

			return e
		}),
		definition.NewOption(), // подключает зависимости и route для работы с бизнес логикой приложения
		fx.Invoke(app.startHTTPServer),
	)

	ctx := context.Background()
	if err := fxApp.Start(ctx); err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	if err := fxApp.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func (app *App) startHTTPServer(lifecycle fx.Lifecycle, e *echo.Echo) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				e.GET("/liveness", func(c echo.Context) error {
					return c.NoContent(http.StatusOK)
				})
				e.GET("/readiness", func(c echo.Context) error {
					//if err := app.DB.Ping(); err != nil {
					//	return c.NoContent(http.StatusServiceUnavailable)
					//}

					return c.NoContent(http.StatusOK)
				})

				err := e.Start(fmt.Sprintf(":%d", app.Config.HTTP().Port))
				if err != nil && errors.Is(err, http.ErrServerClosed) {
					e.Logger.Fatal("server shutdown", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
			defer cancel()

			e.Logger.Info("shutting down server...")

			return e.Shutdown(ctx)
		},
	})
}
