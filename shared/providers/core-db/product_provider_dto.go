package core_db

import (
	"sort"

	"shared/entity"
	"shared/providers/core-db/models"
)

func productToEntity(r models.Product) entity.Product {
	out := entity.Product{
		ID:         r.ID,
		Slug:       r.Slug,
		Name:       r.Name,
		CategoryID: r.CategoryID,
	}
	if r.Category != nil {
		c := categoryToEntity(*r.Category)
		out.Category = &c
	}
	if len(r.Images) > 0 {
		out.Images = imagesToEntities(r.Images)
	}
	if len(r.AttributeValues) > 0 {
		out.Attributes = attributeValuesToEntities(r.AttributeValues)
	}
	return out
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

// attributeValuesToEntities резолвит relation Attribute → entity.AttributeValue
// и сортирует по Attribute.SortOrder.
func attributeValuesToEntities(rows []models.ProductAttributeValue) []entity.AttributeValue {
	out := make([]entity.AttributeValue, 0, len(rows))
	for _, r := range rows {
		if r.Attribute == nil {
			continue
		}
		out = append(out, entity.AttributeValue{
			Attribute: attributeToEntity(*r.Attribute),
			ValueText: r.ValueText,
			ValueInt:  r.ValueInt,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Attribute.SortOrder < out[j].Attribute.SortOrder
	})
	return out
}