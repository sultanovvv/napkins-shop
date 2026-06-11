package cart

import (
	"napkins-shop/public-executor/api"
	"shared/entity"
	"shared/providers/s3storage"
)

// dto переиспользует тот же IURLBuilder, что и productHandler, чтобы
// картинки в позициях корзины строились ровно по той же логике.
type dto struct {
	urls s3storage.IURLBuilder
}

func newDTO(urls s3storage.IURLBuilder) dto { return dto{urls: urls} }

func (d dto) cart(c *entity.Cart) api.CartResponse {
	if c == nil {
		return api.CartResponse{Items: []api.CartItem{}, TotalQuantity: 0}
	}
	items := make([]api.CartItem, len(c.Items))
	total := 0
	for i, it := range c.Items {
		items[i] = d.cartItem(it)
		total += it.Quantity
	}
	return api.CartResponse{Items: items, TotalQuantity: total}
}

func (d dto) cartItem(it entity.CartItem) api.CartItem {
	out := api.CartItem{
		ProductId: it.ProductID,
		Quantity:  it.Quantity,
	}
	if it.Product != nil {
		p := d.product(it.Product)
		out.Product = &p
	}
	return out
}

// product почти точь-в-точь как в product/handler_dto.go — не вытаскиваем в
// shared, чтобы не плодить циклы импорта и не размывать слой. Любая правка
// формата Product должна повторно отразиться в обоих местах (всего две точки).
func (d dto) product(p *entity.Product) api.Product {
	out := api.Product{
		Id:         p.ID,
		Slug:       p.Slug,
		Name:       p.Name,
		Attributes: attrs(p.Attributes),
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

func attrs(values []entity.AttributeValue) []api.ProductAttributeValue {
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