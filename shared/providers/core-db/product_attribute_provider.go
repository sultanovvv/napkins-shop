package core_db

import (
	"context"

	"github.com/uptrace/bun"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

type IProductAttributeProvider interface {
	ListByProduct(ctx context.Context, productID int64) ([]entity.AttributeValueRow, error)
	ListByProducts(ctx context.Context, productIDs []int64) ([]entity.AttributeValueRow, error)
}

type productAttributeProvider struct{ db postgres.IBaseProvider }

func NewProductAttributeProvider(db postgres.IBaseProvider) IProductAttributeProvider {
	return &productAttributeProvider{db: db}
}

func (p *productAttributeProvider) ListByProduct(ctx context.Context, productID int64) ([]entity.AttributeValueRow, error) {
	var rows []models.ProductAttributeValue
	err := p.db.Conn(ctx).NewSelect().Model(&rows).
		Where("product_id = ?", productID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return productAttrToEntities(rows), nil
}

func (p *productAttributeProvider) ListByProducts(ctx context.Context, productIDs []int64) ([]entity.AttributeValueRow, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	var rows []models.ProductAttributeValue
	err := p.db.Conn(ctx).NewSelect().Model(&rows).
		Where("product_id IN (?)", bun.In(productIDs)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return productAttrToEntities(rows), nil
}

func productAttrToEntities(rows []models.ProductAttributeValue) []entity.AttributeValueRow {
	out := make([]entity.AttributeValueRow, len(rows))
	for i, r := range rows {
		out[i] = entity.AttributeValueRow{
			ProductID:   r.ProductID,
			AttributeID: r.AttributeID,
			ValueText:   r.ValueText,
			ValueInt:    r.ValueInt,
		}
	}
	return out
}