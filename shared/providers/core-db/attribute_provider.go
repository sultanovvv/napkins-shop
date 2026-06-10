package core_db

import (
	"context"
	"database/sql"
	"errors"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

type IAttributeProvider interface {
	List(ctx context.Context) ([]entity.Attribute, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Attribute, error)
	GetByID(ctx context.Context, id int64) (*entity.Attribute, error)
}

type attributeProvider struct{ db postgres.IBaseProvider }

func NewAttributeProvider(db postgres.IBaseProvider) IAttributeProvider {
	return &attributeProvider{db: db}
}

func (p *attributeProvider) List(ctx context.Context) ([]entity.Attribute, error) {
	var rows []models.Attribute
	err := p.db.Conn(ctx).NewSelect().Model(&rows).Order("sort_order", "name").Scan(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Attribute, len(rows))
	for i := range rows {
		out[i] = attributeToEntity(rows[i])
	}
	return out, nil
}

func (p *attributeProvider) GetBySlug(ctx context.Context, slug string) (*entity.Attribute, error) {
	row := models.Attribute{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("slug = ?", slug).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := attributeToEntity(row)
	return &e, nil
}

func (p *attributeProvider) GetByID(ctx context.Context, id int64) (*entity.Attribute, error) {
	row := models.Attribute{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := attributeToEntity(row)
	return &e, nil
}

func attributeToEntity(r models.Attribute) entity.Attribute {
	return entity.Attribute{
		ID:        r.ID,
		Slug:      r.Slug,
		Name:      r.Name,
		ValueType: entity.ValueType(r.ValueType),
		SortOrder: r.SortOrder,
	}
}