package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"napkins-shop/public-executor/bootstrap/config"
	"napkins-shop/public-executor/bootstrap/storage"
	"napkins-shop/public-executor/internal/definition"
	"napkins-shop/public-executor/internal/utils/logger"

	core_db "shared/providers/core-db"
)

const shutdownTimeout = 10 * time.Second

type App struct {
	Logger *zap.Logger
	Config *config.AppConfig
}

func NewApp(configPath string) (*App, error) {
	config.InitViperByEnv(configPath)
	cfg := config.NewAppConfig()
	log := logger.NewLogger()

	return &App{
		Logger: log,
		Config: cfg,
	}, nil
}

func (app *App) Start() error {
	fxApp := fx.New(
		fx.Provide(func() *config.AppConfig { return app.Config }),
		fx.Provide(func() *zap.Logger { return app.Logger }),

		// БД: подключение поднимаем через shared, дальше упаковываем в
		// Connection и раздаём IBaseProvider / ITransactionProvider —
		// чтобы новые сервисы (admin) использовали ровно ту же обвязку.
		fx.Provide(func() (*bun.DB, error) {
			return core_db.InitDatabase(*app.Config.Database())
		}),
		fx.Provide(core_db.NewConnection),
		fx.Provide(core_db.NewBaseProvider),
		fx.Provide(core_db.NewTransactionProvider),

		fx.Provide(func() (*minio.Client, error) {
			return storage.NewMinioClient(app.Config)
		}),
		fx.Provide(func() *echo.Echo {
			e := echo.New()
			e.Use(echozap.ZapLogger(app.Logger))
			e.Use(middleware.RequestID())
			e.Use(middleware.CORS())
			e.Use(middleware.Recover())

			return e
		}),
		definition.NewOption(),
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
