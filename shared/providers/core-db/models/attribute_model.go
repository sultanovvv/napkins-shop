package models

import "github.com/uptrace/bun"

// Attribute — bun-модель таблицы attributes. Внутренний тип persistence;
// используется как relation-target в ProductAttributeValue.
type Attribute struct {
	bun.BaseModel `bun:"table:attributes,alias:a"`

	ID        int64  `bun:"id,pk,autoincrement"`
	Slug      string `bun:"slug,unique,notnull"`
	Name      string `bun:"name,notnull"`
	ValueType string `bun:"value_type,notnull"`
	SortOrder int    `bun:"sort_order,notnull,default:0"`
}

// ProductAttributeValue — bun-модель таблицы product_attribute_values.
// Attribute подгружается через .Relation("AttributeValues.Attribute") в
// провайдере; на маппере собирается entity.AttributeValue.
type ProductAttributeValue struct {
	bun.BaseModel `bun:"table:product_attribute_values"`

	ProductID   int64   `bun:"product_id,pk"`
	AttributeID int64   `bun:"attribute_id,pk"`
	ValueText   *string `bun:"value_text"`
	ValueInt    *int    `bun:"value_int"`

	Attribute *Attribute `bun:"rel:belongs-to,join:attribute_id=id"`
}
