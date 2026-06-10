package order

// CreateOrderInUDTO — состав заказа на оформление.
type CreateOrderInUDTO struct {
	Items []CreateOrderItemInUDTO
}

// CreateOrderItemInUDTO — одна позиция в заказе: id товара + количество.
type CreateOrderItemInUDTO struct {
	ProductID int64
	Quantity  int
}
