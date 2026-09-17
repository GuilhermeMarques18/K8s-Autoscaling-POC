package main

import (
	"log"

	server "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/server"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	undo, err := maxprocs.Set(maxprocs.Logger(log.Printf))
	defer undo()
	if err != nil {
		log.Printf("Falha ao configura maxProcess: %v", err)
	}

	httpServer := server.NewHttpServer(":8081")
	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	grpcServer := server.NewGRPCServer(":8080")
	if err := grpcServer.Run(); err != nil {
		log.Fatalf("GRPC server error: %v", err)
	}
}
