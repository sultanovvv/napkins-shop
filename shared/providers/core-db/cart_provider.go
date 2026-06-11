package core_db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

var ErrCartNotFound = errors.New("cart not found")

type ICartProvider interface {
	GetByUserID(ctx context.Context, userID int64) (*entity.Cart, error)
	GetByGuestHash(ctx context.Context, guestHash string) (*entity.Cart, error)

	// EnsureForUser возвращает корзину пользователя, создавая пустую, если нет.
	EnsureForUser(ctx context.Context, userID int64) (*entity.Cart, error)
	// EnsureForGuest возвращает корзину гостя, создавая пустую, если нет.
	EnsureForGuest(ctx context.Context, guestHash string) (*entity.Cart, error)

	// UpsertItem прибавляет qty к существующей позиции либо создаёт новую.
	UpsertItem(ctx context.Context, cartID, productID int64, qty int) error
	// SetItemQuantity жёстко выставляет qty (qty<=0 → удалить позицию).
	SetItemQuantity(ctx context.Context, cartID, productID int64, qty int) error
	// RemoveItem удаляет одну позицию.
	RemoveItem(ctx context.Context, cartID, productID int64) error
	// Clear удаляет все позиции корзины (саму корзину не удаляет).
	Clear(ctx context.Context, cartID int64) error

	// MoveAllItems переносит позиции из srcCart в dstCart, суммируя
	// quantity по совпадающим product_id. srcCart остаётся пустым.
	MoveAllItems(ctx context.Context, srcCartID, dstCartID int64) error
	// Delete удаляет корзину целиком (вместе с cart_items по FK CASCADE).
	Delete(ctx context.Context, cartID int64) error
}

type cartProvider struct{ db postgres.IBaseProvider }

func NewCartProvider(db postgres.IBaseProvider) ICartProvider {
	return &cartProvider{db: db}
}

func (p *cartProvider) GetByUserID(ctx context.Context, userID int64) (*entity.Cart, error) {
	return p.getOne(ctx, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("c.user_id = ?", userID)
	})
}

func (p *cartProvider) GetByGuestHash(ctx context.Context, guestHash string) (*entity.Cart, error) {
	return p.getOne(ctx, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("c.guest_token_hash = ?", guestHash)
	})
}

func (p *cartProvider) EnsureForUser(ctx context.Context, userID int64) (*entity.Cart, error) {
	existing, err := p.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	row := models.Cart{UserID: &userID}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		Returning("id, created_at, updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := cartToEntity(row)
	return &e, nil
}

func (p *cartProvider) EnsureForGuest(ctx context.Context, guestHash string) (*entity.Cart, error) {
	existing, err := p.GetByGuestHash(ctx, guestHash)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	row := models.Cart{GuestTokenHash: &guestHash}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		Returning("id, created_at, updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := cartToEntity(row)
	return &e, nil
}

func (p *cartProvider) UpsertItem(ctx context.Context, cartID, productID int64, qty int) error {
	if qty <= 0 {
		return nil
	}
	row := models.CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  qty,
	}
	_, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		On("CONFLICT (cart_id, product_id) DO UPDATE SET quantity = ci.quantity + EXCLUDED.quantity").
		Exec(ctx)
	if err != nil {
		return err
	}
	return p.touchCart(ctx, cartID)
}

func (p *cartProvider) SetItemQuantity(ctx context.Context, cartID, productID int64, qty int) error {
	if qty <= 0 {
		return p.RemoveItem(ctx, cartID, productID)
	}
	row := models.CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  qty,
	}
	_, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		On("CONFLICT (cart_id, product_id) DO UPDATE SET quantity = EXCLUDED.quantity").
		Exec(ctx)
	if err != nil {
		return err
	}
	return p.touchCart(ctx, cartID)
}

func (p *cartProvider) RemoveItem(ctx context.Context, cartID, productID int64) error {
	_, err := p.db.Conn(ctx).NewDelete().
		Model((*models.CartItem)(nil)).
		Where("cart_id = ?", cartID).
		Where("product_id = ?", productID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return p.touchCart(ctx, cartID)
}

func (p *cartProvider) Clear(ctx context.Context, cartID int64) error {
	_, err := p.db.Conn(ctx).NewDelete().
		Model((*models.CartItem)(nil)).
		Where("cart_id = ?", cartID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return p.touchCart(ctx, cartID)
}

// MoveAllItems сливает srcCart → dstCart одним SQL: INSERT … SELECT с
// ON CONFLICT, суммирующим quantity. После — позиции src удаляются.
// Вызывать ОБЯЗАТЕЛЬНО внутри ITransactionProvider.Run, чтобы dst и src
// были консистентны при сбое.
func (p *cartProvider) MoveAllItems(ctx context.Context, srcCartID, dstCartID int64) error {
	conn := p.db.Conn(ctx)
	if _, err := conn.NewRaw(`
		INSERT INTO cart_items (cart_id, product_id, quantity, added_at)
		SELECT ?, product_id, quantity, added_at
		FROM cart_items
		WHERE cart_id = ?
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, dstCartID, srcCartID).Exec(ctx); err != nil {
		return err
	}
	if _, err := conn.NewDelete().
		Model((*models.CartItem)(nil)).
		Where("cart_id = ?", srcCartID).
		Exec(ctx); err != nil {
		return err
	}
	return p.touchCart(ctx, dstCartID)
}

func (p *cartProvider) Delete(ctx context.Context, cartID int64) error {
	_, err := p.db.Conn(ctx).NewDelete().
		Model((*models.Cart)(nil)).
		Where("id = ?", cartID).
		Exec(ctx)
	return err
}

func (p *cartProvider) touchCart(ctx context.Context, cartID int64) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.Cart)(nil)).
		Set("updated_at = NOW()").
		Where("id = ?", cartID).
		Exec(ctx)
	return err
}

// getOne — общий путь GetByUserID/GetByGuestHash. Подгружает Items.Product,
// чтобы вызывающему было что отрисовать.
func (p *cartProvider) getOne(ctx context.Context, predicate func(*bun.SelectQuery) *bun.SelectQuery) (*entity.Cart, error) {
	row := models.Cart{}
	q := p.db.Conn(ctx).NewSelect().
		Model(&row).
		Relation("Items", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("added_at")
		}).
		Relation("Items.Product")
	q = predicate(q)
	err := q.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := cartToEntity(row)
	return &e, nil
}