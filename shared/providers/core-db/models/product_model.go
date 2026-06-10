package models

import "github.com/uptrace/bun"

// Product — bun-модель таблицы products. Внутренний тип persistence;
// провайдер мапит её в entity.Product на границе.
type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID         int64  `bun:"id,pk,autoincrement"`
	Slug       string `bun:"slug,unique,notnull"`
	Name       string `bun:"name,notnull"`
	CategoryID *int64 `bun:"category_id"`
}