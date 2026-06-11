package models

import "github.com/uptrace/bun"

// Category — bun-модель таблицы categories. Внутренний тип persistence;
// провайдер мапит её в entity.Category на границе.
type Category struct {
	bun.BaseModel `bun:"table:categories,alias:c"`

	ID        int64  `bun:"id,pk,autoincrement"`
	Slug      string `bun:"slug,unique,notnull"`
	Name      string `bun:"name,notnull"`
	ParentID  *int64 `bun:"parent_id"`
	SortOrder int    `bun:"sort_order,notnull,default:0"`
}