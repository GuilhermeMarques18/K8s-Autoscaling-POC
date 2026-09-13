package types

import (
	"context"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *orders.Order) error
	GetOrders(ctx context.Context, customerID int32) ([]*orders.Order, error)
}
