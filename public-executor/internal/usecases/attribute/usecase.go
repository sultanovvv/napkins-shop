package attribute

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"shared/entity"
	core_db "shared/providers/core-db"
)

var ErrNotFound = errors.New("attribute not found")

type IAttributeUseCases interface {
	List(ctx context.Context) ([]entity.Attribute, error)
	GetByID(ctx context.Context, id int64) (*entity.Attribute, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Attribute, error)
}

type useCases struct {
	logger   *zap.Logger
	provider core_db.IAttributeProvider
}

func NewUseCase(provider core_db.IAttributeProvider, logger *zap.Logger) IAttributeUseCases {
	return &useCases{provider: provider, logger: logger}
}

func (u *useCases) List(ctx context.Context) ([]entity.Attribute, error) {
	return u.provider.List(ctx)
}

func (u *useCases) GetByID(ctx context.Context, id int64) (*entity.Attribute, error) {
	a, err := u.provider.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}

func (u *useCases) GetBySlug(ctx context.Context, slug string) (*entity.Attribute, error) {
	a, err := u.provider.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}