package http

import (
	"strings"

	"github.com/labstack/echo/v4"

	"napkins-shop/public-executor/api"
	"napkins-shop/public-executor/internal/middlewares"

	authhandler "napkins-shop/public-executor/internal/entrypoints/http/auth"
	carthandler "napkins-shop/public-executor/internal/entrypoints/http/cart"
	categoryhandler "napkins-shop/public-executor/internal/entrypoints/http/category"
	orderhandler "napkins-shop/public-executor/internal/entrypoints/http/order"
	producthandler "napkins-shop/public-executor/internal/entrypoints/http/product"
)

var _ api.ServerInterface = (*handlerContainer)(nil)

type handlerContainer struct {
	*producthandler.ProductHandler
	*categoryhandler.CategoryHandler
	*orderhandler.OrderHandler
	*authhandler.AuthHandler
	*carthandler.CartHandler
}

// Register монтирует все публичные эндпоинты на echo через сгенерированный
// oapi-codegen ServerInterface, а затем навешивает auth/guest-cart
// middleware'ы селективно через WithSkipper. Селекция по path даёт чистое
// разделение: глобальный e.Use не задевает /openapi.yaml и /swagger, а
// per-route переустановка handler'ов не нужна (всё описано в openapi.yaml).
func Register(
	e *echo.Echo,
	product *producthandler.ProductHandler,
	category *categoryhandler.CategoryHandler,
	order *orderhandler.OrderHandler,
	auth *authhandler.AuthHandler,
	cart *carthandler.CartHandler,
	authMW *middlewares.AuthMiddleware,
	guestCartMW *middlewares.GuestCartMiddleware,
) {
	// OptionalUser применяется ко всем запросам — если Bearer присутствует,
	// в echo-context появится user_id; на путях без авторизации это no-op.
	e.Use(authMW.OptionalUser())

	// GuestCart cookie выдаётся/читается только на путях, где он нужен:
	// корзина и login/register (там его надо смержить).
	e.Use(middlewares.WithSkipper(
		func(c echo.Context) bool {
			p := c.Request().URL.Path
			return !strings.HasPrefix(p, "/api/v1/public/cart/") &&
				p != "/api/v1/public/auth/login" &&
				p != "/api/v1/public/auth/register"
		},
		guestCartMW.Handle(),
	))

	// RequireUser — только на /auth/logoutAll. Расширять по мере появления
	// приватных операций (например, /me, /orders/myOrders).
	e.Use(middlewares.WithSkipper(
		func(c echo.Context) bool {
			return c.Request().URL.Path != "/api/v1/public/auth/logoutAll"
		},
		authMW.RequireUser(),
	))

	api.RegisterHandlers(e, &handlerContainer{
		ProductHandler:  product,
		CategoryHandler: category,
		OrderHandler:    order,
		AuthHandler:     auth,
		CartHandler:     cart,
	})
}