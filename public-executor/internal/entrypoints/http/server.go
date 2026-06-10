package http

import (
	"github.com/labstack/echo/v4"

	"napkins-shop/public-executor/internal/api"
	categoryhandler "napkins-shop/public-executor/internal/entrypoints/http/category"
	orderhandler "napkins-shop/public-executor/internal/entrypoints/http/order"
	producthandler "napkins-shop/public-executor/internal/entrypoints/http/product"
)

var _ api.ServerInterface = (*handlerContainer)(nil)

type handlerContainer struct {
	*producthandler.ProductHandler
	*categoryhandler.CategoryHandler
	*orderhandler.OrderHandler
}

// Register монтирует все публичные эндпоинты на echo через сгенерированный
// oapi-codegen ServerInterface.
func Register(
	e *echo.Echo,
	product *producthandler.ProductHandler,
	category *categoryhandler.CategoryHandler,
	order *orderhandler.OrderHandler,
) {
	api.RegisterHandlers(e, &handlerContainer{
		ProductHandler:  product,
		CategoryHandler: category,
		OrderHandler:    order,
	})
}