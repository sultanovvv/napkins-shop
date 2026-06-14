package cart

// Actor — кто действует с корзиной. Ровно один из UserID/GuestToken не пуст.
// Middleware складывает оба в echo-context; usecase сам выбирает, по чему
// искать (UserID имеет приоритет).
type Actor struct {
	UserID     *int64
	GuestToken string
}

// AddItemInUDTO — добавить (или увеличить qty) для товара.
type AddItemInUDTO struct {
	Actor     Actor
	ProductID int64
	Quantity  int
}

// SetQuantityInUDTO — жёстко выставить количество.
type SetQuantityInUDTO struct {
	Actor     Actor
	ProductID int64
	Quantity  int
}

// RemoveItemInUDTO — удалить позицию полностью.
type RemoveItemInUDTO struct {
	Actor     Actor
	ProductID int64
}