package core_db

import (
	"context"
	"database/sql"
	"errors"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

// ProductListFilter — фильтр для IProductProvider.List.
// CategorySlug пустой = вся витрина; иначе фильтр по slug категории,
// включая всех её потомков (если slug указывает на родителя).
type ProductListFilter struct {
	CategorySlug string
}

type IProductProvider interface {
	List(ctx context.Context, filter ProductListFilter) ([]entity.Product, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Product, error)
	GetByID(ctx context.Context, id int64) (*entity.Product, error)
}

type productProvider struct{ db postgres.IBaseProvider }

func NewProductProvider(db postgres.IBaseProvider) IProductProvider {
	return &productProvider{db: db}
}

func (p *productProvider) List(ctx context.Context, filter ProductListFilter) ([]entity.Product, error) {
	var rows []models.Product
	q := p.db.Conn(ctx).NewSelect().Model(&rows)
	if filter.CategorySlug != "" {
		q = q.Where(
			`category_id IN (
				SELECT id FROM categories
				WHERE slug = ?
				   OR parent_id = (SELECT id FROM categories WHERE slug = ?)
			)`,
			filter.CategorySlug, filter.CategorySlug,
		)
	}
	if err := q.Order("name").Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]entity.Product, len(rows))
	for i := range rows {
		out[i] = productToEntity(rows[i])
	}
	return out, nil
}

func (p *productProvider) GetBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	row := models.Product{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("slug = ?", slug).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := productToEntity(row)
	return &e, nil
}

func (p *productProvider) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	row := models.Product{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := productToEntity(row)
	return &e, nil
}

func productToEntity(r models.Product) entity.Product {
	return entity.Product{
		ID:         r.ID,
		Slug:       r.Slug,
		Name:       r.Name,
		CategoryID: r.CategoryID,
	}
}