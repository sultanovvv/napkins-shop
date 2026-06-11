package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"napkins-shop/public-executor/configs"
	core_db "shared/providers/core-db"
	"shared/providers/s3storage"
	"shared/utils"
)

type App struct {
	db      *bun.DB
	minio   *minio.Client
	appConf *configs.AppConfig
	server  *echo.Echo
	logger  *zap.Logger
}

func NewApp() (*App, error) {
	cfg, err := configs.NewAppConfig()
	if err != nil {
		return nil, err
	}

	db, err := core_db.InitDatabase(*cfg.Postgres)
	if err != nil {
		return nil, err
	}

	minioClient, err := s3storage.NewClient(*cfg.S3)
	if err != nil {
		return nil, err
	}

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	return &App{
		db:      db,
		minio:   minioClient,
		appConf: cfg,
		server:  server,
		logger:  utils.NewLogger(),
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	fxApp := fx.New(
		Infrastructure(a),
		Providers(),
		fx.Invoke(a.start),
	)

	if err := fxApp.Start(ctx); err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	if err := fxApp.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) start(lifecycle fx.Lifecycle, e *echo.Echo) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				e.GET("/liveness", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
				e.GET("/readiness", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

				a.logger.Info(fmt.Sprintf("http-server started on port: %d", a.appConf.HTTP.Port))
				if err := e.Start(fmt.Sprintf(":%d", a.appConf.HTTP.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
					a.logger.Error("http-server failed", zap.Error(err))
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			a.logger.Info(fmt.Sprintf("http-server shutting down with timeout: %s", a.appConf.ShutdownTimeout))

			ctx, cancel := context.WithTimeout(ctx, a.appConf.ShutdownTimeout)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			done := make(chan struct{})

			go func() {
				wg.Wait()
				close(done)
			}()

			go func() {
				defer wg.Done()
				if err := e.Shutdown(ctx); err != nil {
					a.logger.Error("error on stop http connections", zap.Error(err))
				}
			}()

			select {
			case <-ctx.Done():
			case <-done:
			}

			return nil
		},
	})
}