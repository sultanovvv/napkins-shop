package category

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

func NewPostgresRepository(db *bun.DB) ICategoryRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) List() ([]entity.Category, error) {
	var rows []row
	err := r.db.NewSelect().
		Model(&rows).
		OrderExpr("parent_id NULLS FIRST").
		Order("sort_order").
		Order("name").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]entity.Category, len(rows))
	for i := range rows {
		out[i] = rows[i].toEntity()
	}
	return out, nil
}

func (r *postgresRepository) GetBySlug(slug string) (*entity.Category, error) {
	var c row
	err := r.db.NewSelect().Model(&c).Where("slug = ?", slug).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := c.toEntity()
	return &e, nil
}

func (r *postgresRepository) GetByID(id int64) (*entity.Category, error) {
	var c row
	err := r.db.NewSelect().Model(&c).Where("id = ?", id).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := c.toEntity()
	return &e, nil
}
