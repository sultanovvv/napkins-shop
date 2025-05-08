package product

import (
	"go.uber.org/zap"
)

type useCases struct {
	logger *zap.Logger
}

type IProductUseCases interface {
}

func NewUseCase(
	logger *zap.Logger) IProductUseCases {
	return &useCases{
		logger: logger,
	}
}
