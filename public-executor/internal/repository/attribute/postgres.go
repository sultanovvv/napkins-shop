package attribute

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

func NewPostgresRepository(db *bun.DB) IAttributeRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) List() ([]entity.Attribute, error) {
	var rows []row
	err := r.db.NewSelect().Model(&rows).Order("sort_order", "name").Scan(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]entity.Attribute, len(rows))
	for i := range rows {
		out[i] = rows[i].toEntity()
	}
	return out, nil
}

func (r *postgresRepository) GetBySlug(slug string) (*entity.Attribute, error) {
	var a row
	err := r.db.NewSelect().Model(&a).Where("slug = ?", slug).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := a.toEntity()
	return &e, nil
}

func (r *postgresRepository) GetByID(id int64) (*entity.Attribute, error) {
	var a row
	err := r.db.NewSelect().Model(&a).Where("id = ?", id).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := a.toEntity()
	return &e, nil
}
