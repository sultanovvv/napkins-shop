package attribute

import (
	"sync"

	"shared/entity"
)

type mockRepository struct {
	mu   sync.RWMutex
	rows []entity.Attribute
}

func NewMockRepository() IAttributeRepository {
	return &mockRepository{
		rows: []entity.Attribute{
			{ID: 1, Slug: "size", Name: "Размер", ValueType: entity.ValueTypeString, SortOrder: 10},
			{ID: 2, Slug: "layers", Name: "Слои", ValueType: entity.ValueTypeInt, SortOrder: 20},
			{ID: 3, Slug: "color", Name: "Цвет", ValueType: entity.ValueTypeString, SortOrder: 30},
			{ID: 4, Slug: "theme", Name: "Тема", ValueType: entity.ValueTypeString, SortOrder: 40},
			{ID: 5, Slug: "pattern", Name: "Узор", ValueType: entity.ValueTypeString, SortOrder: 50},
		},
	}
}

func (r *mockRepository) List() ([]entity.Attribute, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]entity.Attribute, len(r.rows))
	copy(out, r.rows)
	return out, nil
}

func (r *mockRepository) GetBySlug(slug string) (*entity.Attribute, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.rows {
		if r.rows[i].Slug == slug {
			a := r.rows[i]
			return &a, nil
		}
	}
	return nil, nil
}

func (r *mockRepository) GetByID(id int64) (*entity.Attribute, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.rows {
		if r.rows[i].ID == id {
			a := r.rows[i]
			return &a, nil
		}
	}
	return nil, nil
}
