package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"napkins-shop/backoffice-executor/configs"
	"shared/utils"
)

type App struct {
	appConf *configs.AppConfig
	server  *echo.Echo
	logger  *zap.Logger
}

func NewApp() (*App, error) {
	cfg, err := configs.NewAppConfig()
	if err != nil {
		return nil, err
	}

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Server.ReadHeaderTimeout = 10 * time.Second
	server.Server.IdleTimeout = 60 * time.Second

	return &App{appConf: cfg, server: server, logger: utils.NewLogger()}, nil
}

func (a *App) Start(ctx context.Context) error {
	fxApp := fx.New(Infrastructure(a), Providers(), fx.Invoke(a.start))
	if err := fxApp.Start(ctx); err != nil {
		return err
	}

	var runErr error
	select {
	case <-ctx.Done():
	case signal := <-fxApp.Wait():
		if signal.ExitCode != 0 {
			runErr = fmt.Errorf("application stopped with exit code %d", signal.ExitCode)
		}
	}

	// Остановка использует отдельный контекст, даже если входной уже отменён.
	stopCtx, cancel := context.WithTimeout(context.Background(), a.appConf.ShutdownTimeout)
	defer cancel()
	return errors.Join(runErr, fxApp.Stop(stopCtx))
}

func (a *App) start(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			address := net.JoinHostPort(a.appConf.HTTP.Host, strconv.Itoa(a.appConf.HTTP.Port))
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return fmt.Errorf("listen on %s: %w", address, err)
			}
			a.server.Server.Handler = a.server
			go func() {
				if err := a.server.Server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					a.logger.Error("http server failed", zap.Error(err))
					if err := shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
						a.logger.Error("request application shutdown", zap.Error(err))
					}
				}
			}()
			a.logger.Info("http server started", zap.String("address", address))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			a.logger.Info("http server shutting down")
			if err := a.server.Shutdown(ctx); err != nil {
				return errors.Join(err, a.server.Close())
			}
			return nil
		},
	})
}
