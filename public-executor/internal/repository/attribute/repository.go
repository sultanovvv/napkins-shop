package attribute

import (
	"github.com/uptrace/bun"

	"shared/entity"
)

type row struct {
	bun.BaseModel `bun:"table:attributes,alias:a"`

	ID        int64             `bun:"id,pk,autoincrement"`
	Slug      string            `bun:"slug,unique,notnull"`
	Name      string            `bun:"name,notnull"`
	ValueType entity.ValueType  `bun:"value_type,notnull"`
	SortOrder int               `bun:"sort_order,notnull,default:0"`
}

func (r row) toEntity() entity.Attribute {
	return entity.Attribute{
		ID:        r.ID,
		Slug:      r.Slug,
		Name:      r.Name,
		ValueType: r.ValueType,
		SortOrder: r.SortOrder,
	}
}

type IAttributeRepository interface {
	List() ([]entity.Attribute, error)
	GetBySlug(slug string) (*entity.Attribute, error)
	GetByID(id int64) (*entity.Attribute, error)
}
