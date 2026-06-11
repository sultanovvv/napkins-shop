package category

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"shared/entity"
	core_db "shared/providers/core-db"
)

var ErrNotFound = errors.New("category not found")

type ICategoryUseCases interface {
	GetTree(ctx context.Context) ([]entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	List(ctx context.Context) ([]entity.Category, error)
}

type useCases struct {
	logger   *zap.Logger
	provider core_db.ICategoryProvider
}

func NewUseCase(provider core_db.ICategoryProvider, logger *zap.Logger) ICategoryUseCases {
	return &useCases{provider: provider, logger: logger}
}

func (u *useCases) GetTree(ctx context.Context) ([]entity.Category, error) {
	rows, err := u.provider.List(ctx)
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]*entity.Category, len(rows))
	for i := range rows {
		c := rows[i]
		byID[c.ID] = &c
	}

	var roots []*entity.Category
	for i := range rows {
		node := byID[rows[i].ID]
		if rows[i].ParentID == nil {
			roots = append(roots, node)
			continue
		}
		parent, ok := byID[*rows[i].ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, *node)
	}

	out := make([]entity.Category, len(roots))
	for i, r := range roots {
		out[i] = *r
	}
	return out, nil
}

func (u *useCases) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	c, err := u.provider.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

func (u *useCases) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	c, err := u.provider.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

func (u *useCases) List(ctx context.Context) ([]entity.Category, error) {
	return u.provider.List(ctx)
}