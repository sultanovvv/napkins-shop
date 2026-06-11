package entity

import "time"

// Cart — корзина. Ровно одно из UserID/GuestTokenHash не nil
// (DB-уровневый CHECK + частичные UNIQUE-индексы).
type Cart struct {
	ID              int64
	UserID          *int64
	GuestTokenHash  *string
	Items           []CartItem
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CartItem — позиция в корзине. Product опционально подгружается, когда
// usecase отдаёт корзину наружу: для отображения фронт должен знать имя/slug.
type CartItem struct {
	CartID    int64
	ProductID int64
	Quantity  int
	AddedAt   time.Time
	Product   *Product
}