package order

// CreateOrderInUDTO — состав заказа на оформление.
// UserID опционален: для гостевого checkout он nil, тогда проверка email-
// verified пропускается. У залогиненного пользователя — обязателен.
type CreateOrderInUDTO struct {
	UserID *int64
	Items  []CreateOrderItemInUDTO
}

// CreateOrderItemInUDTO — одна позиция в заказе: id товара + количество.
type CreateOrderItemInUDTO struct {
	ProductID int64
	Quantity  int
}
