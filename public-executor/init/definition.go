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
	authentry "napkins-shop/public-executor/internal/entrypoints/http/auth"
	cartentry "napkins-shop/public-executor/internal/entrypoints/http/cart"
	categoryentry "napkins-shop/public-executor/internal/entrypoints/http/category"
	"napkins-shop/public-executor/internal/entrypoints/http/docs"
	orderentry "napkins-shop/public-executor/internal/entrypoints/http/order"
	productentry "napkins-shop/public-executor/internal/entrypoints/http/product"
	"napkins-shop/public-executor/internal/middlewares"
	authuc "napkins-shop/public-executor/internal/usecases/auth"
	cartuc "napkins-shop/public-executor/internal/usecases/cart"
	categoryuc "napkins-shop/public-executor/internal/usecases/category"
	orderuc "napkins-shop/public-executor/internal/usecases/order"
	productuc "napkins-shop/public-executor/internal/usecases/product"
	productimageuc "napkins-shop/public-executor/internal/usecases/product_image"

	authcfg "shared/configs/auth"
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
		fx.Provide(func() *authcfg.Config { return a.appConf.Auth }),
		fx.Provide(func() authcfg.Argon2Config { return a.appConf.Auth.Argon2 }),
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
			core_db.NewUserProvider,
			core_db.NewAuthIdentityProvider,
			core_db.NewRefreshTokenProvider,
			core_db.NewPasswordResetProvider,
			core_db.NewEmailVerificationProvider,
			core_db.NewCartProvider,

			// shared s3 providers
			s3storage.NewURLBuilder,
			s3storage.NewGateway,

			// auth primitives
			authuc.NewArgon2Hasher,
			authuc.NewHS256TokenSigner,
			authuc.NewLocalIdentityProvider,
			provideIdentityRegistry,

			// usecases
			categoryuc.NewUseCase,
			productuc.NewUseCase,
			productimageuc.NewUploader,
			orderuc.NewUseCase,
			provideCartUseCase,
			authuc.NewUseCase,

			// middleware
			middlewares.NewAuthMiddleware,
			middlewares.NewGuestCartMiddleware,

			// http entrypoints
			productentry.NewProductHandler,
			categoryentry.NewCategoryHandler,
			orderentry.NewOrderHandler,
			authentry.NewAuthHandler,
			cartentry.NewCartHandler,
		),
		fx.Invoke(
			middlewares.Register,
			httpentry.Register,
			docs.Register,
		),
	)
}

// provideIdentityRegistry собирает реестр identity-провайдеров. Сейчас
// в нём только локальный (email+password). Когда появятся Google/Yandex/VK —
// добавятся их конструкторы и здесь же укажутся как зависимости.
func provideIdentityRegistry(local authuc.IIdentityProvider) *authuc.Registry {
	return authuc.NewRegistry(local)
}

// provideCartUseCase — отдельная функция, потому что cart-usecase
// реализует authuc.ICartMerger, и auth-usecase зависит от того же
// инстанса. fx сам разрулит граф по интерфейсу, если зарегистрировать
// возвращаемый тип как оба контракта.
func provideCartUseCase(
	logger *zap.Logger,
	tx core_db.ITransactionProvider,
	cartProv core_db.ICartProvider,
	productProv core_db.IProductProvider,
) (cartuc.ICartUseCases, authuc.ICartMerger) {
	uc := cartuc.NewUseCase(logger, tx, cartProv, productProv)
	return uc, uc.(authuc.ICartMerger)
}