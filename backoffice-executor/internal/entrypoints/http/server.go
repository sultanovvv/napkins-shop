package http

import (
	"github.com/labstack/echo/v4"

	"napkins-shop/backoffice-executor/api"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/category"
	"napkins-shop/backoffice-executor/internal/entrypoints/http/health"
)

const BasePath = "/backoffice/v1"

type handlerContainer struct {
	*health.HealthHandler
	*category.CategoryHandler
}

var _ api.ServerInterface = (*handlerContainer)(nil)

func Register(e *echo.Echo, healthHandler *health.HealthHandler, categoryHandler *category.CategoryHandler) {
	api.RegisterHandlers(e, &handlerContainer{HealthHandler: healthHandler, CategoryHandler: categoryHandler})
}
