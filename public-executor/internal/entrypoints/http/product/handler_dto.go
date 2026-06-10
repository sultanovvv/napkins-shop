package product

import (
	"napkins-shop/public-executor/api"
	"shared/entity"
	"shared/providers/s3storage"
)

type dto struct {
	urls s3storage.IURLBuilder
}

func newDTO(urls s3storage.IURLBuilder) dto { return dto{urls: urls} }

func (d dto) product(p *entity.Product) api.Product {
	out := api.Product{
		Id:         p.ID,
		Slug:       p.Slug,
		Name:       p.Name,
		Attributes: toAPIAttributes(p.Attributes),
		Images:     d.images(p.Images),
	}
	if p.Category != nil {
		out.Category = &api.CategoryRef{
			Id:   p.Category.ID,
			Slug: p.Category.Slug,
			Name: p.Category.Name,
		}
	}
	return out
}

func (d dto) images(items []entity.Image) []api.ProductImage {
	out := make([]api.ProductImage, len(items))
	for i, img := range items {
		out[i] = api.ProductImage{
			Url:       d.urls.URL(img.Key),
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
		}
	}
	return out
}

func toAPIAttributes(values []entity.AttributeValue) []api.ProductAttributeValue {
	out := make([]api.ProductAttributeValue, len(values))
	for i, v := range values {
		out[i] = api.ProductAttributeValue{
			Attribute: api.Attribute{
				Id:        v.Attribute.ID,
				Slug:      v.Attribute.Slug,
				Name:      v.Attribute.Name,
				ValueType: api.AttributeValueType(v.Attribute.ValueType),
				SortOrder: v.Attribute.SortOrder,
			},
			ValueText: v.ValueText,
			ValueInt:  v.ValueInt,
		}
	}
	return out
}
