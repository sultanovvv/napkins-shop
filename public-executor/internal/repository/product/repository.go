package product

import (
	"github.com/uptrace/bun"

	"shared/entity"
)

type row struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID         int64  `bun:"id,pk,autoincrement"`
	Slug       string `bun:"slug,unique,notnull"`
	Name       string `bun:"name,notnull"`
	CategoryID *int64 `bun:"category_id"`
}

func (r row) toEntity() entity.Product {
	return entity.Product{
		ID:         r.ID,
		Slug:       r.Slug,
		Name:       r.Name,
		CategoryID: r.CategoryID,
	}
}

type ListFilter struct {
	CategorySlug string
}

type IProductRepository interface {
	List(filter ListFilter) ([]entity.Product, error)
	GetBySlug(slug string) (*entity.Product, error)
	GetByID(id int64) (*entity.Product, error)
}
