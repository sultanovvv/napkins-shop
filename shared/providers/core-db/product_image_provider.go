package core_db

import (
	"context"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

// IProductImageProvider обслуживает запись картинок товара. Чтение идёт
// через relation Product.Images — отдельных list-методов нет.
type IProductImageProvider interface {
	Insert(ctx context.Context, img *entity.Image) error
	ClearPrimary(ctx context.Context, productID int64) error
}

type productImageProvider struct{ db postgres.IBaseProvider }

func NewProductImageProvider(db postgres.IBaseProvider) IProductImageProvider {
	return &productImageProvider{db: db}
}

func (p *productImageProvider) Insert(ctx context.Context, img *entity.Image) error {
	row := imageFromEntity(*img)
	_, err := p.db.Conn(ctx).NewInsert().Model(&row).Returning("*").Exec(ctx)
	if err != nil {
		return err
	}
	*img = imageToEntity(row)
	return nil
}

func (p *productImageProvider) ClearPrimary(ctx context.Context, productID int64) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.ProductImage)(nil)).
		Set("is_primary = false").
		Where("product_id = ?", productID).
		Where("is_primary = true").
		Exec(ctx)
	return err
}
