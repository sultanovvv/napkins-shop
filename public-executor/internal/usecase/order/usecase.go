package order

import (
	"errors"
	"time"

	"go.uber.org/zap"

	orderrepo "napkins-shop/public-executor/internal/repository/order"
	productrepo "napkins-shop/public-executor/internal/repository/product"
	"shared/entity"
)

var ErrNotFound = errors.New("order not found")
var ErrProductNotFound = errors.New("product not found")

type IOrderUseCases interface {
	CreateOrder(in CreateOrderInUDTO) (*entity.Order, error)
	GetOrder(id int64) (*entity.Order, error)
}

type useCases struct {
	logger      *zap.Logger
	orderRepo   orderrepo.IOrderRepository
	productRepo productrepo.IProductRepository
}

func NewUseCase(
	orderRepo orderrepo.IOrderRepository,
	productRepo productrepo.IProductRepository,
	logger *zap.Logger,
) IOrderUseCases {
	return &useCases{
		logger:      logger,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *useCases) CreateOrder(in CreateOrderInUDTO) (*entity.Order, error) {
	items := make([]entity.OrderItem, 0, len(in.Items))
	for _, it := range in.Items {
		p, err := u.productRepo.GetByID(it.ProductID)
		if err != nil || p == nil {
			return nil, ErrProductNotFound
		}
		items = append(items, entity.OrderItem{
			ProductID: it.ProductID,
			Name:      p.Name,
			Quantity:  it.Quantity,
		})
	}

	return u.orderRepo.Create(entity.Order{
		Items:     items,
		Status:    "pending",
		CreatedAt: time.Now(),
	})
}

func (u *useCases) GetOrder(id int64) (*entity.Order, error) {
	o, err := u.orderRepo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return o, nil
}
