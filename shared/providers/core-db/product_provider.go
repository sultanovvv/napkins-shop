package core_db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

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
	q := p.withRelations(p.db.Conn(ctx).NewSelect().Model(&rows))
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
	err := p.withRelations(p.db.Conn(ctx).NewSelect().Model(&row)).
		Where("p.slug = ?", slug).
		Scan(ctx)
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
	err := p.withRelations(p.db.Conn(ctx).NewSelect().Model(&row)).
		Where("p.id = ?", id).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := productToEntity(row)
	return &e, nil
}

// withRelations подцепляет Category (belongs-to), Images (has-many с order'ом)
// и AttributeValues.Attribute (has-many → belongs-to). Bun делает это двумя
// дополнительными IN-запросами, что заметно дешевле прежних батчей.
func (p *productProvider) withRelations(q *bun.SelectQuery) *bun.SelectQuery {
	return q.
		Relation("Category").
		Relation("Images", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("is_primary DESC", "sort_order", "id")
		}).
		Relation("AttributeValues.Attribute")
}