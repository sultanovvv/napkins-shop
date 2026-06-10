package product_attribute

import (
	"context"

	"github.com/uptrace/bun"

	"shared/entity"
)

type postgresRepository struct {
	db *bun.DB
}

func NewPostgresRepository(db *bun.DB) IProductAttributeRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) ListByProduct(productID int64) ([]entity.AttributeValueRow, error) {
	var rows []row
	err := r.db.NewSelect().Model(&rows).
		Where("product_id = ?", productID).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return toEntities(rows), nil
}

func (r *postgresRepository) ListByProducts(productIDs []int64) ([]entity.AttributeValueRow, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	var rows []row
	err := r.db.NewSelect().Model(&rows).
		Where("product_id IN (?)", bun.In(productIDs)).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return toEntities(rows), nil
}

func toEntities(rows []row) []entity.AttributeValueRow {
	out := make([]entity.AttributeValueRow, len(rows))
	for i := range rows {
		out[i] = rows[i].toEntity()
	}
	return out
}
