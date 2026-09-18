package service

import (
	"context"
	"crypto/sha256"
	"os"
	"strconv"
	"sync"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
)

var (
	createWorkIterations = loadWorkIterations("ORDER_WORK_ITERATIONS_CREATE", 20000)
	getWorkIterations     = loadWorkIterations("ORDER_WORK_ITERATIONS_GET", 5000)
)

func loadWorkIterations(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func simulateWork(seed string, iterations int) {
	h := sha256.New()
	data := []byte(seed)
	for i := 0; i < iterations; i++ {
		h.Write(data)
		data = h.Sum(nil)
	}
}

type OrderService struct {
	mu     sync.RWMutex
	orders map[int32][]*orders.Order
}

func NewOrderService() *OrderService {
	return &OrderService{orders: make(map[int32][]*orders.Order)}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	simulateWork(strconv.Itoa(int(order.OrderID))+strconv.Itoa(int(order.ProductID)), createWorkIterations)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.CustomerID] = append(s.orders[order.CustomerID], order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, customerID int32) ([]*orders.Order, error) {
	simulateWork(strconv.Itoa(int(customerID)), getWorkIterations)

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.orders[customerID], nil
}