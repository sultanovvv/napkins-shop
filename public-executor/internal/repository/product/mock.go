package product

import (
	"sync"

	"shared/entity"
)

type mockRepository struct {
	mu       sync.RWMutex
	products []entity.Product
}

func NewMockRepository() IProductRepository {
	int64p := func(v int64) *int64 { return &v }
	return &mockRepository{
		products: []entity.Product{
			{ID: 1, Slug: "roses-napkin", Name: "Салфетка с розами", CategoryID: int64p(8)},
			{ID: 2, Slug: "daisies-napkin", Name: "Салфетка с ромашками", CategoryID: int64p(8)},
			{ID: 3, Slug: "butterflies-napkin", Name: "Салфетка с бабочками", CategoryID: int64p(7)},
			{ID: 4, Slug: "leaves-napkin", Name: "Салфетка с листьями", CategoryID: int64p(10)},
			{ID: 5, Slug: "gnomes-napkin", Name: "Гномы новогодние", CategoryID: int64p(13)},
		},
	}
}

func (r *mockRepository) List(filter ListFilter) ([]entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if filter.CategorySlug == "" {
		out := make([]entity.Product, len(r.products))
		copy(out, r.products)
		return out, nil
	}
	return nil, nil
}

func (r *mockRepository) GetBySlug(slug string) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.products {
		if r.products[i].Slug == slug {
			p := r.products[i]
			return &p, nil
		}
	}
	return nil, nil
}

func (r *mockRepository) GetByID(id int64) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.products {
		if r.products[i].ID == id {
			p := r.products[i]
			return &p, nil
		}
	}
	return nil, nil
}
