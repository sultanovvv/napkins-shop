package core_db

import (
	"context"
	"errors"
	"sync"

	"shared/entity"
)

var ErrOrderNotFound = errors.New("order not found")

type IOrderProvider interface {
	Create(ctx context.Context, order entity.Order) (*entity.Order, error)
	GetByID(ctx context.Context, id int64) (*entity.Order, error)
}

type orderProvider struct {
	mu     sync.RWMutex
	orders map[int64]entity.Order
	nextID int64
}

// NewOrderProvider — пока in-memory: таблицы orders/order_items нет в миграциях.
// Когда появится — заменить тело на postgres-реализацию через postgres.IBaseProvider,
// сохранив сигнатуру.
func NewOrderProvider() IOrderProvider {
	return &orderProvider{
		orders: make(map[int64]entity.Order),
		nextID: 1,
	}
}

func (p *orderProvider) Create(_ context.Context, order entity.Order) (*entity.Order, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	order.ID = p.nextID
	p.nextID++
	p.orders[order.ID] = order
	return &order, nil
}

func (p *orderProvider) GetByID(_ context.Context, id int64) (*entity.Order, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	o, ok := p.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return &o, nil
}