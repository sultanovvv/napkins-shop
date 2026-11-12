package middlewares

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Register(e *echo.Echo, logger *zap.Logger) {
	e.Use(middleware.RequestID())
	e.Use(echozap.ZapLogger(logger))
	e.Use(middleware.Recover())
}
