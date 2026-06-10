package order

import (
	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"napkins-shop/public-executor/internal/api"
	orderusecase "napkins-shop/public-executor/internal/usecase/order"
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

	order, err := h.orderUseCase.CreateOrder(toCreateOrderInput(body))
	if err != nil {
		if errors.Is(err, orderusecase.ErrProductNotFound) {
			return ctx.JSON(400, api.ErrorModel{
				Error: api.ErrorDetailsModel{Code: "PRODUCT_NOT_FOUND", Message: err.Error()},
			})
		}
		h.logger.Error("CreateOrder failed", zap.Error(err))
		return ctx.JSON(500, api.ErrorModel{
			Error: api.ErrorDetailsModel{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
	}

	return ctx.JSON(200, toOrderResponse(order))
}

func (h OrderHandler) GetOrder(ctx echo.Context, id int64) error {
	order, err := h.orderUseCase.GetOrder(id)
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
