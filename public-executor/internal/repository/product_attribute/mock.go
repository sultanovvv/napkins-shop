package product_attribute

import (
	"sync"

	"shared/entity"
)

type mockRepository struct {
	mu   sync.RWMutex
	rows []entity.AttributeValueRow
}

func NewMockRepository() IProductAttributeRepository {
	strp := func(s string) *string { return &s }
	intp := func(v int) *int { return &v }
	return &mockRepository{
		rows: []entity.AttributeValueRow{
			{ProductID: 1, AttributeID: 1, ValueText: strp("33x33 см")},
			{ProductID: 1, AttributeID: 2, ValueInt: intp(3)},
			{ProductID: 1, AttributeID: 3, ValueText: strp("белый")},
			{ProductID: 1, AttributeID: 5, ValueText: strp("розы")},
		},
	}
}

func (r *mockRepository) ListByProduct(productID int64) ([]entity.AttributeValueRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []entity.AttributeValueRow
	for _, v := range r.rows {
		if v.ProductID == productID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (r *mockRepository) ListByProducts(productIDs []int64) ([]entity.AttributeValueRow, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	want := make(map[int64]struct{}, len(productIDs))
	for _, id := range productIDs {
		want[id] = struct{}{}
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []entity.AttributeValueRow
	for _, v := range r.rows {
		if _, ok := want[v.ProductID]; ok {
			out = append(out, v)
		}
	}
	return out, nil
}
