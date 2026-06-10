package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"napkins-shop/public-executor/configs"
	httpentry "napkins-shop/public-executor/internal/entrypoints/http"
	categoryentry "napkins-shop/public-executor/internal/entrypoints/http/category"
	orderentry "napkins-shop/public-executor/internal/entrypoints/http/order"
	productentry "napkins-shop/public-executor/internal/entrypoints/http/product"
	"napkins-shop/public-executor/internal/middlewares"
	categoryuc "napkins-shop/public-executor/internal/usecases/category"
	orderuc "napkins-shop/public-executor/internal/usecases/order"
	productuc "napkins-shop/public-executor/internal/usecases/product"
	productimageuc "napkins-shop/public-executor/internal/usecases/product_image"

	"shared/configs/s3"
	core_db "shared/providers/core-db"
	"shared/providers/s3storage"
)

func Infrastructure(a *App) fx.Option {
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
		fx.Provide(func() *configs.AppConfig { return a.appConf }),
		fx.Provide(func() *zap.Logger { return a.logger }),
		fx.Provide(func() *echo.Echo { return a.server }),
		fx.Provide(func() *minio.Client { return a.minio }),
		fx.Provide(func() *bun.DB { return a.db }),
		fx.Provide(func() *core_db.Connection { return core_db.NewConnection(a.db) }),
		fx.Provide(func() s3.Config { return *a.appConf.S3 }),
		fx.Provide(func() *http.Client { return &http.Client{} }),
	)
}

func Providers() fx.Option {
	return fx.Options(
		fx.Provide(
			// shared core-db providers
			core_db.NewBaseProvider,
			core_db.NewTransactionProvider,
			core_db.NewProductProvider,
			core_db.NewCategoryProvider,
			core_db.NewProductImageProvider,
			core_db.NewOrderProvider,

			// shared s3 providers
			s3storage.NewURLBuilder,
			s3storage.NewGateway,

			// usecases
			categoryuc.NewUseCase,
			productuc.NewUseCase,
			productimageuc.NewUploader,
			orderuc.NewUseCase,

			// http entrypoints
			productentry.NewProductHandler,
			categoryentry.NewCategoryHandler,
			orderentry.NewOrderHandler,
		),
		fx.Invoke(
			middlewares.Register,
			httpentry.Register,
		),
	)
}