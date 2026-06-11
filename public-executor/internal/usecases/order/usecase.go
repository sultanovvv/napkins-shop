package order

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"shared/entity"
	core_db "shared/providers/core-db"
)

var ErrNotFound = errors.New("order not found")
var ErrProductNotFound = errors.New("product not found")
var ErrEmailNotVerified = errors.New("email not verified")

type IOrderUseCases interface {
	CreateOrder(ctx context.Context, in CreateOrderInUDTO) (*entity.Order, error)
	GetOrder(ctx context.Context, id int64) (*entity.Order, error)
}

type useCases struct {
	logger      *zap.Logger
	orderProv   core_db.IOrderProvider
	productProv core_db.IProductProvider
	userProv    core_db.IUserProvider
}

func NewUseCase(
	orderProv core_db.IOrderProvider,
	productProv core_db.IProductProvider,
	userProv core_db.IUserProvider,
	logger *zap.Logger,
) IOrderUseCases {
	return &useCases{
		logger:      logger,
		orderProv:   orderProv,
		productProv: productProv,
		userProv:    userProv,
	}
}

func (u *useCases) CreateOrder(ctx context.Context, in CreateOrderInUDTO) (*entity.Order, error) {
	// Гость (UserID == nil) checkout-ит без ограничений. Залогиненный
	// пользователь должен иметь подтверждённый email.
	if in.UserID != nil {
		user, err := u.userProv.GetByID(ctx, *in.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil || user.EmailVerifiedAt == nil {
			return nil, ErrEmailNotVerified
		}
	}

	items := make([]entity.OrderItem, 0, len(in.Items))
	for _, it := range in.Items {
		p, err := u.productProv.GetByID(ctx, it.ProductID)
		if err != nil || p == nil {
			return nil, ErrProductNotFound
		}
		items = append(items, entity.OrderItem{
			ProductID: it.ProductID,
			Name:      p.Name,
			Quantity:  it.Quantity,
		})
	}

	return u.orderProv.Create(ctx, entity.Order{
		Items:     items,
		Status:    "pending",
		CreatedAt: time.Now(),
	})
}

func (u *useCases) GetOrder(ctx context.Context, id int64) (*entity.Order, error) {
	o, err := u.orderProv.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return o, nil
}