package product_image

import (
	"time"

	"github.com/uptrace/bun"

	"shared/entity"
)

type row struct {
	bun.BaseModel `bun:"table:product_images,alias:pi"`

	ID        int64     `bun:"id,pk,autoincrement"`
	ProductID int64     `bun:"product_id,notnull"`
	S3Key     string    `bun:"s3_key,notnull"`
	SortOrder int       `bun:"sort_order,notnull,default:0"`
	IsPrimary bool      `bun:"is_primary,notnull,default:false"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
}

func (r row) toEntity() entity.Image {
	return entity.Image{
		ID:        r.ID,
		ProductID: r.ProductID,
		Key:       r.S3Key,
		SortOrder: r.SortOrder,
		IsPrimary: r.IsPrimary,
		CreatedAt: r.CreatedAt,
	}
}

func fromEntity(e entity.Image) row {
	return row{
		ID:        e.ID,
		ProductID: e.ProductID,
		S3Key:     e.Key,
		SortOrder: e.SortOrder,
		IsPrimary: e.IsPrimary,
		CreatedAt: e.CreatedAt,
	}
}

type IProductImageRepository interface {
	ListByProduct(productID int64) ([]entity.Image, error)
	ListByProducts(productIDs []int64) ([]entity.Image, error)
	Insert(img *entity.Image) error
	ClearPrimary(productID int64) error
}
