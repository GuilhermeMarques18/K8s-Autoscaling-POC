package handler

import (
	"context"

	orders "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/types"
	"google.golang.org/grpc"
)

type OrdersGrpcHandler struct {
	orders.UnimplementedOrderServiceServer
	ordersService types.OrderService
}

func NewGrpcOrdersService(grpcServer *grpc.Server, ordersService types.OrderService) {
	h := &OrdersGrpcHandler{ordersService: ordersService}
	orders.RegisterOrderServiceServer(grpcServer, h)
}

func (h *OrdersGrpcHandler) GetOrder(ctx context.Context, req *orders.GetOrderRequest) (*orders.GetOrderResponse, error) {
	orderList, err := h.ordersService.GetOrders(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	return &orders.GetOrderResponse{Orders: orderList}, nil
}

func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		OrderID:    42,
		CustomerID: req.CustomerId,
		ProductID:  req.ProductId,
		Quantity:   req.Quantity,
	}

	if err := h.ordersService.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return &orders.CreateOrderResponse{Status: "success"}, nil
}
