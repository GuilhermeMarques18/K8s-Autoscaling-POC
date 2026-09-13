package handler

import (
	"net/http"

	orders "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/util"
	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/types"
)

type OrdersHttpHandler struct {
	ordersService types.OrderService
}

func NewHttpOrdersHandler(ordersService types.OrderService) *OrdersHttpHandler {
	handler := &OrdersHttpHandler{ordersService: ordersService}
	return handler
}

func (h *OrdersHttpHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /orders", h.CreateOrder)
}

func (h *OrdersHttpHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest
	err := util.ParseJSON(r, &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	order := &orders.Order{
		OrderID:    42,
		CustomerID: req.GetCustomerId(),
		ProductID:  req.GetProductId(),
		Quantity:   req.GetQuantity(),
	}

	err = h.ordersService.CreateOrder(r.Context(), order)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	res := &orders.CreateOrderResponse{Status: "success"}
	util.WriteJSON(w, http.StatusOK, res)
}
