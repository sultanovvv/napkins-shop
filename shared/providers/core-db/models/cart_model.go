package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Cart — bun-модель таблицы carts. Одна корзина — либо UserID, либо
// GuestTokenHash; DB-уровневый CHECK гарантирует "ровно один".
// Items подгружаются через .Relation("Items") в провайдере.
type Cart struct {
	bun.BaseModel `bun:"table:carts,alias:c"`

	ID             int64      `bun:"id,pk,autoincrement"`
	UserID         *int64     `bun:"user_id"`
	GuestTokenHash *string    `bun:"guest_token_hash"`
	CreatedAt      time.Time  `bun:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull,default:now()"`

	Items []CartItem `bun:"rel:has-many,join:id=cart_id"`
}

// CartItem — bun-модель таблицы cart_items. Product подгружается через
// .Relation("Items.Product"), когда нужно отдать товар наружу.
type CartItem struct {
	bun.BaseModel `bun:"table:cart_items,alias:ci"`

	CartID    int64     `bun:"cart_id,pk"`
	ProductID int64     `bun:"product_id,pk"`
	Quantity  int       `bun:"quantity,notnull"`
	AddedAt   time.Time `bun:"added_at,notnull,default:now()"`

	Product *Product `bun:"rel:belongs-to,join:product_id=id"`
}