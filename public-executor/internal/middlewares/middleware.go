package middlewares

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Register(e *echo.Echo, log *zap.Logger) error {
	e.Use(echozap.ZapLogger(log))
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	return nil
}