package http

import (
	"github.com/labstack/echo/v4"

	"napkins-shop/public-executor/internal/api"
	categoryhandler "napkins-shop/public-executor/internal/handler/http/public/category"
	orderhandler "napkins-shop/public-executor/internal/handler/http/public/order"
	producthandler "napkins-shop/public-executor/internal/handler/http/public/product"
)

var _ api.ServerInterface = (*handlerContainer)(nil)

type handlerContainer struct {
	*producthandler.ProductHandler
	*categoryhandler.CategoryHandler
	*orderhandler.OrderHandler
}

func RegisterHandlers(
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
