package order

import (
	"napkins-shop/public-executor/api"
	orderusecase "napkins-shop/public-executor/internal/usecases/order"
	"shared/entity"
)

func toCreateOrderInput(body api.CreateOrderRequest) orderusecase.CreateOrderInUDTO {
	items := make([]orderusecase.CreateOrderItemInUDTO, len(body.Items))
	for i, it := range body.Items {
		items[i] = orderusecase.CreateOrderItemInUDTO{
			ProductID: it.ProductId,
			Quantity:  it.Quantity,
		}
	}
	return orderusecase.CreateOrderInUDTO{Items: items}
}

func toOrderResponse(o *entity.Order) api.OrderResponse {
	items := make([]api.OrderItemResponse, len(o.Items))
	for i, it := range o.Items {
		items[i] = api.OrderItemResponse{
			ProductId: it.ProductID,
			Name:      it.Name,
			Quantity:  it.Quantity,
		}
	}
	return api.OrderResponse{
		Id:        o.ID,
		Items:     items,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
	}
}
