package http

import (
	"github.com/labstack/echo/v4"

	"napkins-shop/backoffice-executor/internal/entrypoints/http/health"
)

const BasePath = "/backoffice/v1"

func Register(e *echo.Echo, healthHandler *health.HealthHandler) {
	e.GET(BasePath+"/health", healthHandler.GetHealth)
}
