package entity

import "time"

type Order struct {
	ID        int64
	Items     []OrderItem
	Status    string
	CreatedAt time.Time
}

type OrderItem struct {
	ProductID int64
	Name      string
	Quantity  int
}
