package order

import "shared/entity"

type IOrderRepository interface {
	Create(order entity.Order) (*entity.Order, error)
	GetByID(id int64) (*entity.Order, error)
}
