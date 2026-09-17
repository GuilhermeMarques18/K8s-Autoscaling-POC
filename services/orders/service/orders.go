package service

import (
	"context"
	"crypto/sha256"
	"os"
	"strconv"
	"sync"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
)

var workIterations = loadWorkIterations()

func loadWorkIterations() int {
	if v := os.Getenv("ORDER_WORK_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 20000
}

func simulateProcessing(order *orders.Order) {
	h := sha256.New()
	data := []byte(strconv.Itoa(int(order.OrderID)) + strconv.Itoa(int(order.ProductID)))
	for i := 0; i < workIterations; i++ {
		h.Write(data)
		data = h.Sum(nil)
	}
}

type OrderService struct {
	mu     sync.RWMutex
	orders []*orders.Order
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	simulateProcessing(order)

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
