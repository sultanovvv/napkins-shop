package product

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"shared/entity"
)

type postgresRepository struct {
	db *bun.DB
}

func NewPostgresRepository(db *bun.DB) IProductRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) List(filter ListFilter) ([]entity.Product, error) {
	var rows []row
	q := r.db.NewSelect().Model(&rows)
	if filter.CategorySlug != "" {
		// Покрываем оба уровня дерева: товары категории-листа ИЛИ любого
		// её потомка, если slug указывает на родителя.
		q = q.Where(
			`category_id IN (
				SELECT id FROM categories
				WHERE slug = ?
				   OR parent_id = (SELECT id FROM categories WHERE slug = ?)
			)`,
			filter.CategorySlug, filter.CategorySlug,
		)
	}
	if err := q.Order("name").Scan(context.Background()); err != nil {
		return nil, err
	}
	out := make([]entity.Product, len(rows))
	for i := range rows {
		out[i] = rows[i].toEntity()
	}
	return out, nil
}

func (r *postgresRepository) GetBySlug(slug string) (*entity.Product, error) {
	var p row
	err := r.db.NewSelect().Model(&p).Where("slug = ?", slug).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := p.toEntity()
	return &e, nil
}

func (r *postgresRepository) GetByID(id int64) (*entity.Product, error) {
	var p row
	err := r.db.NewSelect().Model(&p).Where("id = ?", id).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := p.toEntity()
	return &e, nil
}
