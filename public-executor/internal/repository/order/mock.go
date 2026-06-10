package order

import (
	"fmt"
	"sync"

	"shared/entity"
)

type mockRepository struct {
	mu     sync.RWMutex
	orders map[int64]entity.Order
	nextID int64
}

func NewMockRepository() IOrderRepository {
	return &mockRepository{
		orders: make(map[int64]entity.Order),
		nextID: 1,
	}
}

func (r *mockRepository) Create(order entity.Order) (*entity.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = r.nextID
	r.nextID++
	r.orders[order.ID] = order
	return &order, nil
}

func (r *mockRepository) GetByID(id int64) (*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %d not found", id)
	}
	return &o, nil
}
