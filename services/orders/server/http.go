package serve

import (
	"log"
	"net/http"

	handler "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/handler/orders"
	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/service"
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{addr: addr}
}

func (s *HttpServer) Run() error {
	router := http.NewServeMux()

	orderService := service.NewOrderService()
	orderHandler := handler.NewHttpOrdersHandler(orderService)
	orderHandler.RegisterRoutes(router)

	log.Println("Starting HTTP server at ", s.addr)

	return http.ListenAndServe(s.addr, router)

}
