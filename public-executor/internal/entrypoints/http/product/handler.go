package product

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"napkins-shop/public-executor/internal/api"
	"napkins-shop/public-executor/internal/storage/s3url"
	"napkins-shop/public-executor/internal/usecases/product"
	productimageuc "napkins-shop/public-executor/internal/usecases/product_image"
)

const maxUploadBytes = 10 << 20 // 10 MiB

type ProductHandler struct {
	productUseCase product.IProductUseCases
	uploader       productimageuc.IUploader
	dto            dto
	logger         *zap.Logger
}

func NewProductHandler(
	productUseCase product.IProductUseCases,
	uploader productimageuc.IUploader,
	urls s3url.Builder,
	logger *zap.Logger,
) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		uploader:       uploader,
		dto:            newDTO(urls),
		logger:         logger,
	}
}

func (h ProductHandler) GetProduct(ctx echo.Context, slug string) error {
	p, err := h.productUseCase.GetProductBySlug(slug)
	if err != nil {
		if errors.Is(err, product.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, api.ErrorModel{
				Error: api.ErrorDetailsModel{
					Code:    "NOT_FOUND",
					Message: "product not found",
					Details: slug,
				},
			})
		}
		h.logger.Error("GetProduct failed", zap.Error(err))
		return internalErr(ctx)
	}

	return ctx.JSON(http.StatusOK, h.dto.product(p))
}

func (h ProductHandler) GetProductsList(ctx echo.Context, params api.GetProductsListParams) error {
	in := product.GetProductsListInUDTO{}
	if params.Category != nil {
		in.CategorySlug = *params.Category
	}

	items, err := h.productUseCase.GetProductsList(in)
	if err != nil {
		h.logger.Error("GetProductsList failed", zap.Error(err))
		return internalErr(ctx)
	}

	apiItems := make([]api.Product, len(items))
	for i := range items {
		apiItems[i] = h.dto.product(&items[i])
	}

	return ctx.JSON(http.StatusOK, api.GetProductListResponse{
		Items:      apiItems,
		TotalCount: len(apiItems),
	})
}

func (h ProductHandler) UploadProductImage(ctx echo.Context) error {
	productSlug := ctx.FormValue("productSlug")
	if productSlug == "" {
		return badRequest(ctx, "productSlug is required")
	}

	fh, err := ctx.FormFile("file")
	if err != nil {
		return badRequest(ctx, "file is required: "+err.Error())
	}
	if fh.Size <= 0 {
		return badRequest(ctx, "file is empty")
	}
	if fh.Size > maxUploadBytes {
		return badRequest(ctx, "file too large")
	}

	src, err := fh.Open()
	if err != nil {
		h.logger.Error("open multipart file", zap.Error(err))
		return internalErr(ctx)
	}
	defer src.Close()

	isPrimary, _ := parseBool(ctx.FormValue("isPrimary"))
	sortOrder, _ := parseInt(ctx.FormValue("sortOrder"))

	result, err := h.uploader.Upload(ctx.Request().Context(), productimageuc.UploadInUDTO{
		ProductSlug: productSlug,
		Filename:    fh.Filename,
		ContentType: fh.Header.Get("Content-Type"),
		Body:        src,
		Size:        fh.Size,
		IsPrimary:   isPrimary,
		SortOrder:   sortOrder,
	})
	if err != nil {
		switch {
		case errors.Is(err, productimageuc.ErrProductNotFound):
			return ctx.JSON(http.StatusNotFound, api.ErrorModel{
				Error: api.ErrorDetailsModel{Code: "NOT_FOUND", Message: "product not found", Details: productSlug},
			})
		case errors.Is(err, productimageuc.ErrInvalidImage):
			return badRequest(ctx, "invalid image")
		default:
			h.logger.Error("upload image failed", zap.Error(err))
			return internalErr(ctx)
		}
	}

	return ctx.JSON(http.StatusOK, api.UploadProductImageResponse{
		Id:        result.ID,
		Key:       result.Key,
		Url:       result.URL,
		IsPrimary: result.IsPrimary,
		SortOrder: result.SortOrder,
	})
}

func badRequest(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusBadRequest, api.ErrorModel{
		Error: api.ErrorDetailsModel{Code: "BAD_REQUEST", Message: msg},
	})
}

func internalErr(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, api.ErrorModel{
		Error: api.ErrorDetailsModel{Code: "INTERNAL_ERROR", Message: "internal server error"},
	})
}

func parseBool(s string) (bool, bool) {
	switch s {
	case "true", "1", "on", "yes":
		return true, true
	case "false", "0", "off", "no", "":
		return false, true
	default:
		return false, false
	}
}

func parseInt(s string) (int, bool) {
	if s == "" {
		return 0, true
	}
	var v int
	for i, c := range s {
		if c < '0' || c > '9' {
			if i == 0 && c == '-' {
				continue
			}
			return 0, false
		}
		v = v*10 + int(c-'0')
	}
	return v, true
}
