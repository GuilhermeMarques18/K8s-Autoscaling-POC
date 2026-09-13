package handler

import (
	"context"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/types"
)

type OrdersGrpcHandler struct {
	ordersService types.OrderService
	orders.UnimplementedOrderServiceServer
}

func NewGrpcOrdersHandler() {
	gRPCHandler := &OrdersGrpcHandler{}
}

func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		OrderID:    42,
		CustomerId: 2,
		ProductID:  1,
		Quantity:   10,
	}
	err := h.ordersService.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	res := &orders.CreateOrderResponse{
		Status: "sucess",
	}
	return res, nil
}
