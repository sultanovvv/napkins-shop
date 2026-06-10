package product_image

import (
	"context"

	"github.com/uptrace/bun"

	"shared/entity"
)

type postgresRepository struct {
	db *bun.DB
}

func NewPostgresRepository(db *bun.DB) IProductImageRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) ListByProduct(productID int64) ([]entity.Image, error) {
	var rows []row
	err := r.db.NewSelect().Model(&rows).
		Where("product_id = ?", productID).
		Order("is_primary DESC", "sort_order", "id").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return toEntities(rows), nil
}

func (r *postgresRepository) ListByProducts(productIDs []int64) ([]entity.Image, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	var rows []row
	err := r.db.NewSelect().Model(&rows).
		Where("product_id IN (?)", bun.In(productIDs)).
		Order("product_id", "is_primary DESC", "sort_order", "id").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return toEntities(rows), nil
}

func (r *postgresRepository) Insert(img *entity.Image) error {
	rec := fromEntity(*img)
	_, err := r.db.NewInsert().Model(&rec).Returning("*").Exec(context.Background())
	if err != nil {
		return err
	}
	*img = rec.toEntity()
	return nil
}

func (r *postgresRepository) ClearPrimary(productID int64) error {
	_, err := r.db.NewUpdate().
		Model((*row)(nil)).
		Set("is_primary = false").
		Where("product_id = ?", productID).
		Where("is_primary = true").
		Exec(context.Background())
	return err
}

func toEntities(rows []row) []entity.Image {
	out := make([]entity.Image, len(rows))
	for i := range rows {
		out[i] = rows[i].toEntity()
	}
	return out
}
