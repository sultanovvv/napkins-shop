package models

import "github.com/uptrace/bun"

// Product — bun-модель таблицы products. Внутренний тип persistence;
// провайдер мапит её в entity.Product на границе. Relations подгружаются
// через .Relation("...") в провайдере одним SELECT с JOIN/IN.
type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID         int64  `bun:"id,pk,autoincrement"`
	Slug       string `bun:"slug,unique,notnull"`
	Name       string `bun:"name,notnull"`
	CategoryID *int64 `bun:"category_id"`

	Category        *Category               `bun:"rel:belongs-to,join:category_id=id"`
	Images          []ProductImage          `bun:"rel:has-many,join:id=product_id"`
	AttributeValues []ProductAttributeValue `bun:"rel:has-many,join:id=product_id"`
}
