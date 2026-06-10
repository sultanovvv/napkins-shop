package category

import (
	"github.com/uptrace/bun"

	"shared/entity"
)

// row — bun-модель, приватная для репозитория. Наружу отдаём entity.Category.
type row struct {
	bun.BaseModel `bun:"table:categories,alias:c"`

	ID        int64  `bun:"id,pk,autoincrement"`
	Slug      string `bun:"slug,unique,notnull"`
	Name      string `bun:"name,notnull"`
	ParentID  *int64 `bun:"parent_id"`
	SortOrder int    `bun:"sort_order,notnull,default:0"`
}

func (r row) toEntity() entity.Category {
	return entity.Category{
		ID:        r.ID,
		Slug:      r.Slug,
		Name:      r.Name,
		ParentID:  r.ParentID,
		SortOrder: r.SortOrder,
	}
}

type ICategoryRepository interface {
	List() ([]entity.Category, error)
	GetBySlug(slug string) (*entity.Category, error)
	GetByID(id int64) (*entity.Category, error)
}
