package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func categoryToEntity(r models.Category) entity.Category {
	return entity.Category{
		ID:        r.ID,
		Slug:      r.Slug,
		Name:      r.Name,
		ParentID:  r.ParentID,
		SortOrder: r.SortOrder,
	}
}