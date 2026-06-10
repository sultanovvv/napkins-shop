package category

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"napkins-shop/public-executor/api"

	categoryusecase "napkins-shop/public-executor/internal/usecases/category"
)

type CategoryHandler struct {
	categoryUseCase categoryusecase.ICategoryUseCases
	logger          *zap.Logger
}

func NewCategoryHandler(categoryUseCase categoryusecase.ICategoryUseCases, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
		logger:          logger,
	}
}

func (h CategoryHandler) GetCategoriesTree(ctx echo.Context) error {
	tree, err := h.categoryUseCase.GetTree(ctx.Request().Context())
	if err != nil {
		h.logger.Error("GetCategoriesTree failed", zap.Error(err))
		return ctx.JSON(500, api.ErrorModel{
			Error: api.ErrorDetailsModel{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
	}

	items := make([]api.CategoryNode, len(tree))
	for i, c := range tree {
		items[i] = toAPICategoryNode(c)
	}

	return ctx.JSON(200, api.GetCategoriesTreeResponse{Items: items})
}
