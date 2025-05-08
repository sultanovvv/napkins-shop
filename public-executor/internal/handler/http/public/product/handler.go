package product

import (
	_ "errors"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"napkins-shop/public-executor/internal/api"
	"napkins-shop/public-executor/internal/usecase/product"
	_ "net/http"

	_ "github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

type ProductHandler struct {
	productUseCase product.IProductUseCases
	logger         *zap.Logger
}

func NewProductHandler(productUseCase product.IProductUseCases, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
	}
}

func (h ProductHandler) GetProductsList(ctx echo.Context) error {
	out := api.GetProductListResponse{
		Items: []api.Product{{
			Attributes: nil,
			Id:         uuid.New(),
			Name:       "Test product",
		}},
		TotalCount: 1,
	}
	return ctx.JSON(200, out)
}
