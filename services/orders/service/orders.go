package service

import (
	"context"
	"sync"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
)

type OrderService struct {
	mu     sync.Mutex
	orders []*orders.Order
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders = append(s.orders, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, customerID int32) ([]*orders.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []*orders.Order
	for _, o := range s.orders {
		if o.CustomerID == customerID {
			result = append(result, o)
		}
	}
	return result, nil
}
