package serve

import (
	"log"
	"net"

	handler "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/handler/orders"
	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/service"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *GRPCServer {
	return &GRPCServer{addr: addr}
}

func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	orderService := service.NewOrderService()
	handler.NewGrpcOrdersService(grpcServer, orderService)

	log.Println("Starting gRPC server at ", s.addr)

	return grpcServer.Serve(lis)
}
