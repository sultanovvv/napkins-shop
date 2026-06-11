package order

import (
	"errors"
	"napkins-shop/public-executor/api"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"napkins-shop/public-executor/internal/middlewares"
	orderusecase "napkins-shop/public-executor/internal/usecases/order"
)

type OrderHandler struct {
	orderUseCase orderusecase.IOrderUseCases
	logger       *zap.Logger
}

func NewOrderHandler(orderUseCase orderusecase.IOrderUseCases, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

func (h OrderHandler) CreateOrder(ctx echo.Context) error {
	var body api.CreateOrderRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(400, api.ErrorModel{
			Error: api.ErrorDetailsModel{Code: "BAD_REQUEST", Message: "invalid request body"},
		})
	}
	if len(body.Items) == 0 {
		return ctx.JSON(400, api.ErrorModel{
			Error: api.ErrorDetailsModel{Code: "BAD_REQUEST", Message: "items is empty"},
		})
	}

	in := toCreateOrderInput(body)
	// Если запрос пришёл от залогиненного пользователя, OptionalUser
	// middleware положил его user_id в echo-context. Прокидываем — usecase
	// сам решит, проверять ли email_verified_at.
	if uid, ok := ctx.Get(middlewares.CtxUserID).(int64); ok {
		in.UserID = &uid
	}

	order, err := h.orderUseCase.CreateOrder(ctx.Request().Context(), in)
	if err != nil {
		switch {
		case errors.Is(err, orderusecase.ErrProductNotFound):
			return ctx.JSON(400, api.ErrorModel{
				Error: api.ErrorDetailsModel{Code: "PRODUCT_NOT_FOUND", Message: err.Error()},
			})
		case errors.Is(err, orderusecase.ErrEmailNotVerified):
			return ctx.JSON(403, api.ErrorModel{
				Error: api.ErrorDetailsModel{Code: "EMAIL_NOT_VERIFIED", Message: "confirm your email before placing an order"},
			})
		default:
			h.logger.Error("CreateOrder failed", zap.Error(err))
			return ctx.JSON(500, api.ErrorModel{
				Error: api.ErrorDetailsModel{Code: "INTERNAL_ERROR", Message: "internal server error"},
			})
		}
	}

	return ctx.JSON(200, toOrderResponse(order))
}

func (h OrderHandler) GetOrder(ctx echo.Context, id int64) error {
	order, err := h.orderUseCase.GetOrder(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, orderusecase.ErrNotFound) {
			return ctx.JSON(404, api.ErrorModel{
				Error: api.ErrorDetailsModel{Code: "NOT_FOUND", Message: "order not found"},
			})
		}
		h.logger.Error("GetOrder failed", zap.Error(err))
		return ctx.JSON(500, api.ErrorModel{
			Error: api.ErrorDetailsModel{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
	}

	return ctx.JSON(200, toOrderResponse(order))
}
