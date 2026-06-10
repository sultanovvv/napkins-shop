package core_db

import (
	"context"
	"database/sql"
	"errors"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

type ICategoryProvider interface {
	List(ctx context.Context) ([]entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
}

type categoryProvider struct{ db postgres.IBaseProvider }

func NewCategoryProvider(db postgres.IBaseProvider) ICategoryProvider {
	return &categoryProvider{db: db}
}

func (p *categoryProvider) List(ctx context.Context) ([]entity.Category, error) {
	var rows []models.Category
	err := p.db.Conn(ctx).NewSelect().
		Model(&rows).
		OrderExpr("parent_id NULLS FIRST").
		Order("sort_order").
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Category, len(rows))
	for i := range rows {
		out[i] = categoryToEntity(rows[i])
	}
	return out, nil
}

func (p *categoryProvider) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	row := models.Category{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("slug = ?", slug).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := categoryToEntity(row)
	return &e, nil
}

func (p *categoryProvider) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	row := models.Category{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := categoryToEntity(row)
	return &e, nil
}

func categoryToEntity(r models.Category) entity.Category {
	return entity.Category{
		ID:        r.ID,
		Slug:      r.Slug,
		Name:      r.Name,
		ParentID:  r.ParentID,
		SortOrder: r.SortOrder,
	}
}