package http

import (
	"github.com/labstack/echo/v4"
	"napkins-shop/public-executor/internal/api"
	"napkins-shop/public-executor/internal/handler/http/public/product"
)

var _ api.ServerInterface = (*handlerContainer)(nil)

type handlerContainer struct {
	*product.ProductHandler
}

func RegisterHandlers(
	e *echo.Echo,
	product *product.ProductHandler,
) {
	api.RegisterHandlers(e, product)
}
