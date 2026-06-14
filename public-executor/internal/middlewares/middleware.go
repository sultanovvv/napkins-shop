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
	// CORS под credentials: дефолтный `middleware.CORS()` отдаёт "*",
	// и браузер блокирует cookies/Authorization вместе. Дев-фронт ходит с
	// http://localhost:3000 (Next.js), прод-домены добавим, когда появятся.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowCredentials: true,
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderAuthorization,
			echo.HeaderAccept,
			echo.HeaderXRequestedWith,
		},
		MaxAge: 600,
	}))
	e.Use(middleware.Recover())
	return nil
}

// WithSkipper оборачивает middleware так, чтобы оно срабатывало только
// для запросов, на которых skipper возвращает false. Полезно, когда
// `e.Use(...)` нужно ограничить набором роутов (например, RequireUser —
// только на /auth/logoutAll).
func WithSkipper(skip func(c echo.Context) bool, mw echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skip(c) {
				return next(c)
			}
			return mw(next)(c)
		}
	}
}
