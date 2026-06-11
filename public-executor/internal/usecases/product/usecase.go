package product

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"shared/entity"
	core_db "shared/providers/core-db"
)

var ErrNotFound = errors.New("product not found")

type IProductUseCases interface {
	GetProductsList(ctx context.Context, in GetProductsListInUDTO) ([]entity.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*entity.Product, error)
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)
}

type useCases struct {
	logger      *zap.Logger
	productProv core_db.IProductProvider
}

func NewUseCase(productProv core_db.IProductProvider, logger *zap.Logger) IProductUseCases {
	return &useCases{productProv: productProv, logger: logger}
}

func (u *useCases) GetProductsList(ctx context.Context, in GetProductsListInUDTO) ([]entity.Product, error) {
	return u.productProv.List(ctx, core_db.ProductListFilter{CategorySlug: in.CategorySlug})
}

func (u *useCases) GetProductBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	p, err := u.productProv.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrNotFound
	}
	return p, nil
}

func (u *useCases) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	p, err := u.productProv.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrNotFound
	}
	return p, nil
}