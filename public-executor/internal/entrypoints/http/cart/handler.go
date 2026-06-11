package cart

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"napkins-shop/public-executor/api"
	"napkins-shop/public-executor/internal/middlewares"
	cartuc "napkins-shop/public-executor/internal/usecases/cart"
	"shared/providers/s3storage"
)

type CartHandler struct {
	useCase cartuc.ICartUseCases
	dto     dto
	logger  *zap.Logger
}

func NewCartHandler(useCase cartuc.ICartUseCases, urls s3storage.IURLBuilder, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		useCase: useCase,
		dto:     newDTO(urls),
		logger:  logger,
	}
}

func (h *CartHandler) GetCart(ctx echo.Context) error {
	actor := actorFromCtx(ctx)
	cart, err := h.useCase.GetCart(ctx.Request().Context(), actor)
	if err != nil {
		return h.mapErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, h.dto.cart(cart))
}

func (h *CartHandler) CartAddItem(ctx echo.Context) error {
	var body api.CartAddItemRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	cart, err := h.useCase.AddItem(ctx.Request().Context(), cartuc.AddItemInUDTO{
		Actor:     actorFromCtx(ctx),
		ProductID: body.ProductId,
		Quantity:  body.Quantity,
	})
	if err != nil {
		return h.mapErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, h.dto.cart(cart))
}

func (h *CartHandler) CartSetQuantity(ctx echo.Context) error {
	var body api.CartSetQuantityRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	cart, err := h.useCase.SetQuantity(ctx.Request().Context(), cartuc.SetQuantityInUDTO{
		Actor:     actorFromCtx(ctx),
		ProductID: body.ProductId,
		Quantity:  body.Quantity,
	})
	if err != nil {
		return h.mapErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, h.dto.cart(cart))
}

func (h *CartHandler) CartRemoveItem(ctx echo.Context) error {
	var body api.CartRemoveItemRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	cart, err := h.useCase.RemoveItem(ctx.Request().Context(), cartuc.RemoveItemInUDTO{
		Actor:     actorFromCtx(ctx),
		ProductID: body.ProductId,
	})
	if err != nil {
		return h.mapErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, h.dto.cart(cart))
}

func (h *CartHandler) CartClear(ctx echo.Context) error {
	cart, err := h.useCase.Clear(ctx.Request().Context(), actorFromCtx(ctx))
	if err != nil {
		return h.mapErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, h.dto.cart(cart))
}

// actorFromCtx собирает actor'а из echo-context: middleware-auth кладёт
// user_id, middleware-guest-cart кладёт guest-token. user_id побеждает
// (см. cart.useCases.ensure).
func actorFromCtx(ctx echo.Context) cartuc.Actor {
	out := cartuc.Actor{}
	if uid, ok := ctx.Get(middlewares.CtxUserID).(int64); ok {
		out.UserID = &uid
	}
	if guest, ok := ctx.Get(middlewares.CtxGuestToken).(string); ok && guest != "" {
		out.GuestToken = guest
	}
	return out
}

func (h *CartHandler) mapErr(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, cartuc.ErrActorRequired):
		return errorJSON(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "actor missing")
	case errors.Is(err, cartuc.ErrProductNotFound):
		return errorJSON(ctx, http.StatusNotFound, "PRODUCT_NOT_FOUND", "product not found")
	case errors.Is(err, cartuc.ErrInvalidQuantity):
		return badRequest(ctx, "quantity must be non-negative")
	default:
		h.logger.Error("cart op failed", zap.Error(err))
		return internalErr(ctx)
	}
}

func badRequest(ctx echo.Context, msg string) error {
	return errorJSON(ctx, http.StatusBadRequest, "BAD_REQUEST", msg)
}

func internalErr(ctx echo.Context) error {
	return errorJSON(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func errorJSON(ctx echo.Context, status int, code, msg string) error {
	return ctx.JSON(status, api.ErrorModel{
		Error: api.ErrorDetailsModel{
			Code:    code,
			Message: msg,
		},
	})
}