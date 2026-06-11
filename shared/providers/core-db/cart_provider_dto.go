package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func cartToEntity(r models.Cart) entity.Cart {
	out := entity.Cart{
		ID:             r.ID,
		UserID:         r.UserID,
		GuestTokenHash: r.GuestTokenHash,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	if len(r.Items) > 0 {
		out.Items = cartItemsToEntities(r.Items)
	}
	return out
}

func cartItemsToEntities(rows []models.CartItem) []entity.CartItem {
	out := make([]entity.CartItem, len(rows))
	for i, r := range rows {
		item := entity.CartItem{
			CartID:    r.CartID,
			ProductID: r.ProductID,
			Quantity:  r.Quantity,
			AddedAt:   r.AddedAt,
		}
		if r.Product != nil {
			p := productToEntity(*r.Product)
			item.Product = &p
		}
		out[i] = item
	}
	return out
}