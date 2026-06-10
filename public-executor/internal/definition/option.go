package definition

import (
	"go.uber.org/fx"

	httproute "napkins-shop/public-executor/internal/definition/http"
	categoryhandler "napkins-shop/public-executor/internal/handler/http/public/category"
	orderhandler "napkins-shop/public-executor/internal/handler/http/public/order"
	producthandler "napkins-shop/public-executor/internal/handler/http/public/product"
	attributerepo "napkins-shop/public-executor/internal/repository/attribute"
	categoryrepo "napkins-shop/public-executor/internal/repository/category"
	orderrepo "napkins-shop/public-executor/internal/repository/order"
	productrepo "napkins-shop/public-executor/internal/repository/product"
	productattrrepo "napkins-shop/public-executor/internal/repository/product_attribute"
	productimagerepo "napkins-shop/public-executor/internal/repository/product_image"
	"napkins-shop/public-executor/internal/storage/s3gateway"
	"napkins-shop/public-executor/internal/storage/s3url"
	attributeusecase "napkins-shop/public-executor/internal/usecase/attribute"
	categoryusecase "napkins-shop/public-executor/internal/usecase/category"
	orderusecase "napkins-shop/public-executor/internal/usecase/order"
	productusecase "napkins-shop/public-executor/internal/usecase/product"
	productimageusecase "napkins-shop/public-executor/internal/usecase/product_image"
)

// NewOption — fx-граф зависимостей. Postgres-репозитории работают; order пока
// in-memory, пока не появится отдельная таблица заказов.
func NewOption() fx.Option {
	return fx.Options(
		// client + storage
		fx.Provide(httproute.ProvideDefaultHTTPClient),
		fx.Provide(s3url.NewBuilder),
		fx.Provide(s3gateway.NewGateway),

		// repositories
		fx.Provide(productrepo.NewPostgresRepository),
		fx.Provide(categoryrepo.NewPostgresRepository),
		fx.Provide(attributerepo.NewPostgresRepository),
		fx.Provide(productattrrepo.NewPostgresRepository),
		fx.Provide(productimagerepo.NewPostgresRepository),
		fx.Provide(func() orderrepo.IOrderRepository { return orderrepo.NewMockRepository() }),

		// usecases
		fx.Provide(categoryusecase.NewUseCase),
		fx.Provide(attributeusecase.NewUseCase),
		fx.Provide(productusecase.NewUseCase),
		fx.Provide(productimageusecase.NewUploader),
		fx.Provide(orderusecase.NewUseCase),

		// http handlers
		fx.Provide(producthandler.NewProductHandler),
		fx.Provide(categoryhandler.NewCategoryHandler),
		fx.Provide(orderhandler.NewOrderHandler),

		// routes
		fx.Invoke(httproute.RegisterHandlers),
	)
}
