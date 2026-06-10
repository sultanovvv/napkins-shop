package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

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
