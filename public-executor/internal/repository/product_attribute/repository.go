package product_attribute

import (
	"github.com/uptrace/bun"

	"shared/entity"
)

// row — сырая строка product_attribute_values. Маппится в entity.AttributeValueRow.
type row struct {
	bun.BaseModel `bun:"table:product_attribute_values"`

	ProductID   int64   `bun:"product_id,pk"`
	AttributeID int64   `bun:"attribute_id,pk"`
	ValueText   *string `bun:"value_text"`
	ValueInt    *int    `bun:"value_int"`
}

func (r row) toEntity() entity.AttributeValueRow {
	return entity.AttributeValueRow{
		ProductID:   r.ProductID,
		AttributeID: r.AttributeID,
		ValueText:   r.ValueText,
		ValueInt:    r.ValueInt,
	}
}

type IProductAttributeRepository interface {
	ListByProduct(productID int64) ([]entity.AttributeValueRow, error)
	ListByProducts(productIDs []int64) ([]entity.AttributeValueRow, error)
}
