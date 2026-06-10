package models

import (
	"time"

	"github.com/uptrace/bun"
)

// ProductImage — bun-модель таблицы product_images. Внутренний тип
// persistence; провайдер мапит её в entity.Image на границе.
type ProductImage struct {
	bun.BaseModel `bun:"table:product_images,alias:pi"`

	ID        int64     `bun:"id,pk,autoincrement"`
	ProductID int64     `bun:"product_id,notnull"`
	S3Key     string    `bun:"s3_key,notnull"`
	SortOrder int       `bun:"sort_order,notnull,default:0"`
	IsPrimary bool      `bun:"is_primary,notnull,default:false"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
}