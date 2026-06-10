package category

import (
	"sync"

	"shared/entity"
)

type mockRepository struct {
	mu   sync.RWMutex
	rows []entity.Category
}

func NewMockRepository() ICategoryRepository {
	int64p := func(v int64) *int64 { return &v }
	return &mockRepository{
		rows: []entity.Category{
			{ID: 1, Slug: "priroda", Name: "Природа", SortOrder: 10},
			{ID: 2, Slug: "prazdniki-grp", Name: "Праздники", SortOrder: 20},
			{ID: 3, Slug: "tematika", Name: "Тематика", SortOrder: 30},
			{ID: 4, Slug: "gorod-byt", Name: "Город и быт", SortOrder: 40},
			{ID: 5, Slug: "eda-napitki", Name: "Еда и напитки", SortOrder: 50},
			{ID: 6, Slug: "pticy", Name: "Птицы", ParentID: int64p(1), SortOrder: 10},
			{ID: 7, Slug: "zhivotnyy-mir", Name: "Животный мир", ParentID: int64p(1), SortOrder: 20},
			{ID: 8, Slug: "cvety", Name: "Цветы", ParentID: int64p(1), SortOrder: 30},
			{ID: 9, Slug: "venochki", Name: "Веночки разные", ParentID: int64p(1), SortOrder: 40},
			{ID: 10, Slug: "uzory-ornamenty", Name: "Узоры. Орнаменты. Фоны.", ParentID: int64p(1), SortOrder: 50},
			{ID: 11, Slug: "novyy-god", Name: "Новый год и Рождество", ParentID: int64p(2), SortOrder: 10},
			{ID: 12, Slug: "prazdniki", Name: "Праздники", ParentID: int64p(2), SortOrder: 20},
			{ID: 13, Slug: "gnomy", Name: "Гномы", ParentID: int64p(2), SortOrder: 30},
		},
	}
}

func (r *mockRepository) List() ([]entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]entity.Category, len(r.rows))
	copy(out, r.rows)
	return out, nil
}

func (r *mockRepository) GetBySlug(slug string) (*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.rows {
		if r.rows[i].Slug == slug {
			c := r.rows[i]
			return &c, nil
		}
	}
	return nil, nil
}

func (r *mockRepository) GetByID(id int64) (*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.rows {
		if r.rows[i].ID == id {
			c := r.rows[i]
			return &c, nil
		}
	}
	return nil, nil
}
