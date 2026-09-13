package server

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	orders "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/common/genproto/orders"
)

type HttpServer struct {
	addr   string
	client orders.OrderServiceClient
}

func NewHttpServer(httpAddr string, ordersGrpcAddr string) *HttpServer {
	conn := NewGRPCClient(ordersGrpcAddr)
	client := orders.NewOrderServiceClient(conn)

	return &HttpServer{addr: httpAddr, client: client}
}

func (s *HttpServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /", s.listOrders)
	router.HandleFunc("POST /orders", s.createOrder)

	log.Println("Starting HTTP server at ", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *HttpServer) listOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	var customerID int32
	if raw := r.URL.Query().Get("customerId"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			http.Error(w, "invalid customerId", http.StatusBadRequest)
			return
		}
		customerID = int32(parsed)
	}

	res, err := s.client.GetOrder(ctx, &orders.GetOrderRequest{CustomerID: customerID})
	if err != nil {
		log.Printf("client error: %v", err)
		http.Error(w, "failed to fetch orders", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("orders").Parse(ordersTemplate)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, res.Orders); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

type createOrderRequest struct {
	CustomerId int32 `json:"customerId"`
	ProductId  int32 `json:"productId"`
	Quantity   int32 `json:"quantity"`
}

func (s *HttpServer) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	_, err := s.client.CreateOrder(ctx, &orders.CreateOrderRequest{
		CustomerId: req.CustomerId,
		ProductId:  req.ProductId,
		Quantity:   req.Quantity,
	})
	if err != nil {
		log.Printf("client error: %v", err)
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

var ordersTemplate = `<!DOCTYPE html>
<html>
<head>
<title>Kitchen Orders</title>
</head>
<body>
<h1>Orders List</h1>
<table border="1">
<tr>
<th>Order ID</th>
<th>Customer ID</th>
<th>Quantity</th>
</tr>
{{range .}}
<tr>
<td>{{.OrderID}}</td>
<td>{{.CustomerID}}</td>
<td>{{.Quantity}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`
