package definition

import (
	"go.uber.org/fx"
	httproute "napkins-shop/public-executor/internal/definition/http"
	producthandler "napkins-shop/public-executor/internal/handler/http/public/product"
	productusecase "napkins-shop/public-executor/internal/usecase/product"
)

// NewOption определяет зависимости и вызывает route
func NewOption() fx.Option {
	return fx.Options(
		// client
		fx.Provide(
			httproute.ProvideDefaultHTTPClient,
			//httproute.ProvideS3Client,
		),
		// usecase
		fx.Provide(productusecase.NewUseCase),
		// http  handler
		fx.Provide(producthandler.NewProductHandler),
		// route
		fx.Invoke(
			httproute.RegisterHandlers,
		),
	)
}
