package cart

import (
	"context"
	"errors"

	"go.uber.org/zap"

	authuc "napkins-shop/public-executor/internal/usecases/auth"
	"shared/entity"
	core_db "shared/providers/core-db"
)

var (
	ErrActorRequired   = errors.New("cart actor is required (user or guest)")
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidQuantity = errors.New("quantity must be positive")
)

type ICartUseCases interface {
	GetCart(ctx context.Context, actor Actor) (*entity.Cart, error)
	AddItem(ctx context.Context, in AddItemInUDTO) (*entity.Cart, error)
	SetQuantity(ctx context.Context, in SetQuantityInUDTO) (*entity.Cart, error)
	RemoveItem(ctx context.Context, in RemoveItemInUDTO) (*entity.Cart, error)
	Clear(ctx context.Context, actor Actor) (*entity.Cart, error)

	// MergeGuestIntoUser — точка входа для auth-usecase после успешного логина.
	// Гостевая корзина (если найдена) переносится в пользовательскую и удаляется.
	MergeGuestIntoUser(ctx context.Context, guestToken string, userID int64) error
}

type useCases struct {
	logger      *zap.Logger
	tx          core_db.ITransactionProvider
	cartProv    core_db.ICartProvider
	productProv core_db.IProductProvider
}

func NewUseCase(
	logger *zap.Logger,
	tx core_db.ITransactionProvider,
	cartProv core_db.ICartProvider,
	productProv core_db.IProductProvider,
) ICartUseCases {
	return &useCases{
		logger:      logger,
		tx:          tx,
		cartProv:    cartProv,
		productProv: productProv,
	}
}

// Compile-time проверка: useCases удовлетворяет авторовскому ICartMerger.
var _ authuc.ICartMerger = (*useCases)(nil)

func (u *useCases) GetCart(ctx context.Context, actor Actor) (*entity.Cart, error) {
	// ensure создаёт пустую корзину, если её ещё нет, чтобы Get-маршрут не
	// падал на первой загрузке страницы; relations подтянем отдельным SELECT.
	if _, err := u.ensure(ctx, actor); err != nil {
		return nil, err
	}
	return u.readCart(ctx, actor)
}

func (u *useCases) AddItem(ctx context.Context, in AddItemInUDTO) (*entity.Cart, error) {
	if in.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if err := u.assertProductExists(ctx, in.ProductID); err != nil {
		return nil, err
	}

	var resultActor = in.Actor
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		cart, err := u.ensure(ctx, in.Actor)
		if err != nil {
			return err
		}
		return u.cartProv.UpsertItem(ctx, cart.ID, in.ProductID, in.Quantity)
	})
	if err != nil {
		return nil, err
	}
	return u.readCart(ctx, resultActor)
}

func (u *useCases) SetQuantity(ctx context.Context, in SetQuantityInUDTO) (*entity.Cart, error) {
	if in.Quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if in.Quantity > 0 {
		if err := u.assertProductExists(ctx, in.ProductID); err != nil {
			return nil, err
		}
	}

	err := u.tx.Run(ctx, func(ctx context.Context) error {
		cart, err := u.ensure(ctx, in.Actor)
		if err != nil {
			return err
		}
		return u.cartProv.SetItemQuantity(ctx, cart.ID, in.ProductID, in.Quantity)
	})
	if err != nil {
		return nil, err
	}
	return u.readCart(ctx, in.Actor)
}

func (u *useCases) RemoveItem(ctx context.Context, in RemoveItemInUDTO) (*entity.Cart, error) {
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		cart, err := u.ensure(ctx, in.Actor)
		if err != nil {
			return err
		}
		return u.cartProv.RemoveItem(ctx, cart.ID, in.ProductID)
	})
	if err != nil {
		return nil, err
	}
	return u.readCart(ctx, in.Actor)
}

func (u *useCases) Clear(ctx context.Context, actor Actor) (*entity.Cart, error) {
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		cart, err := u.ensure(ctx, actor)
		if err != nil {
			return err
		}
		return u.cartProv.Clear(ctx, cart.ID)
	})
	if err != nil {
		return nil, err
	}
	return u.readCart(ctx, actor)
}

func (u *useCases) MergeGuestIntoUser(ctx context.Context, guestToken string, userID int64) error {
	if guestToken == "" {
		return nil
	}
	guestHash := authuc.HashSecret(guestToken)

	guest, err := u.cartProv.GetByGuestHash(ctx, guestHash)
	if err != nil {
		return err
	}
	if guest == nil {
		return nil
	}

	userCart, err := u.cartProv.EnsureForUser(ctx, userID)
	if err != nil {
		return err
	}

	if err := u.cartProv.MoveAllItems(ctx, guest.ID, userCart.ID); err != nil {
		return err
	}
	return u.cartProv.Delete(ctx, guest.ID)
}

// ensure возвращает (либо создаёт) корзину под актора без подгрузки items —
// дешевле для путей add/remove/set.
func (u *useCases) ensure(ctx context.Context, actor Actor) (*entity.Cart, error) {
	if actor.UserID != nil {
		return u.cartProv.EnsureForUser(ctx, *actor.UserID)
	}
	if actor.GuestToken != "" {
		return u.cartProv.EnsureForGuest(ctx, authuc.HashSecret(actor.GuestToken))
	}
	return nil, ErrActorRequired
}

// readCart возвращает корзину с Items.Product, для отдачи наружу.
func (u *useCases) readCart(ctx context.Context, actor Actor) (*entity.Cart, error) {
	if actor.UserID != nil {
		c, err := u.cartProv.GetByUserID(ctx, *actor.UserID)
		if err != nil {
			return nil, err
		}
		return ensureNotNil(c, *actor.UserID, nil), nil
	}
	if actor.GuestToken != "" {
		hash := authuc.HashSecret(actor.GuestToken)
		c, err := u.cartProv.GetByGuestHash(ctx, hash)
		if err != nil {
			return nil, err
		}
		return ensureNotNil(c, 0, &hash), nil
	}
	return nil, ErrActorRequired
}

// ensureNotNil — если корзины внезапно не оказалось (например, Clear вынес
// все позиции, и кто-то параллельно удалил пустую), возвращаем "пустышку"
// с правильным owner'ом, чтобы handler мог отдать пустой список фронту.
func ensureNotNil(c *entity.Cart, userID int64, guestHash *string) *entity.Cart {
	if c != nil {
		return c
	}
	out := &entity.Cart{}
	if userID != 0 {
		out.UserID = &userID
	}
	if guestHash != nil {
		out.GuestTokenHash = guestHash
	}
	return out
}

func (u *useCases) assertProductExists(ctx context.Context, productID int64) error {
	p, err := u.productProv.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrProductNotFound
	}
	return nil
}