package core_db

import (
	"context"

	"github.com/uptrace/bun"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

type IProductImageProvider interface {
	ListByProduct(ctx context.Context, productID int64) ([]entity.Image, error)
	ListByProducts(ctx context.Context, productIDs []int64) ([]entity.Image, error)
	Insert(ctx context.Context, img *entity.Image) error
	ClearPrimary(ctx context.Context, productID int64) error
}

type productImageProvider struct{ db postgres.IBaseProvider }

func NewProductImageProvider(db postgres.IBaseProvider) IProductImageProvider {
	return &productImageProvider{db: db}
}

func (p *productImageProvider) ListByProduct(ctx context.Context, productID int64) ([]entity.Image, error) {
	var rows []models.ProductImage
	err := p.db.Conn(ctx).NewSelect().Model(&rows).
		Where("product_id = ?", productID).
		Order("is_primary DESC", "sort_order", "id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return imagesToEntities(rows), nil
}

func (p *productImageProvider) ListByProducts(ctx context.Context, productIDs []int64) ([]entity.Image, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	var rows []models.ProductImage
	err := p.db.Conn(ctx).NewSelect().Model(&rows).
		Where("product_id IN (?)", bun.In(productIDs)).
		Order("product_id", "is_primary DESC", "sort_order", "id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return imagesToEntities(rows), nil
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

func imageToEntity(r models.ProductImage) entity.Image {
	return entity.Image{
		ID:        r.ID,
		ProductID: r.ProductID,
		Key:       r.S3Key,
		SortOrder: r.SortOrder,
		IsPrimary: r.IsPrimary,
		CreatedAt: r.CreatedAt,
	}
}

func imageFromEntity(e entity.Image) models.ProductImage {
	return models.ProductImage{
		ID:        e.ID,
		ProductID: e.ProductID,
		S3Key:     e.Key,
		SortOrder: e.SortOrder,
		IsPrimary: e.IsPrimary,
		CreatedAt: e.CreatedAt,
	}
}

func imagesToEntities(rows []models.ProductImage) []entity.Image {
	out := make([]entity.Image, len(rows))
	for i := range rows {
		out[i] = imageToEntity(rows[i])
	}
	return out
}