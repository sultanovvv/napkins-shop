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
	attributerepo "napkins-shop/public-executor/internal/repository/attribute"
	categoryrepo "napkins-shop/public-executor/internal/repository/category"
	orderrepo "napkins-shop/public-executor/internal/repository/order"
	productrepo "napkins-shop/public-executor/internal/repository/product"
	productattrrepo "napkins-shop/public-executor/internal/repository/product_attribute"
	productimagerepo "napkins-shop/public-executor/internal/repository/product_image"
	"napkins-shop/public-executor/internal/storage/s3gateway"
	"napkins-shop/public-executor/internal/storage/s3url"
	attributeuc "napkins-shop/public-executor/internal/usecases/attribute"
	categoryuc "napkins-shop/public-executor/internal/usecases/category"
	orderuc "napkins-shop/public-executor/internal/usecases/order"
	productuc "napkins-shop/public-executor/internal/usecases/product"
	productimageuc "napkins-shop/public-executor/internal/usecases/product_image"

	core_db "shared/providers/core-db"
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
		fx.Provide(func() *http.Client { return &http.Client{} }),
	)
}

func Providers() fx.Option {
	return fx.Options(
		fx.Provide(
			// shared providers (core-db)
			core_db.NewBaseProvider,
			core_db.NewTransactionProvider,

			// storage
			s3url.NewBuilder,
			s3gateway.NewGateway,

			// repositories
			productrepo.NewPostgresRepository,
			categoryrepo.NewPostgresRepository,
			attributerepo.NewPostgresRepository,
			productattrrepo.NewPostgresRepository,
			productimagerepo.NewPostgresRepository,
			func() orderrepo.IOrderRepository { return orderrepo.NewMockRepository() },

			// usecases
			categoryuc.NewUseCase,
			attributeuc.NewUseCase,
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