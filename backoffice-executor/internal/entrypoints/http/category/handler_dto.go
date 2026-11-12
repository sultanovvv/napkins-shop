package category

import (
	"napkins-shop/backoffice-executor/api"
	"shared/entity"
)

func toAPI(item entity.Category) api.Category {
	return api.Category{
		Id: item.ID, Name: item.Name, Slug: item.Slug,
		ParentId: item.ParentID, SortOrder: item.SortOrder,
	}
}
