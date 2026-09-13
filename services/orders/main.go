package main

import (
	"log"

	server "github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/orders/server"
)

func main() {
	httpServer := server.NewHttpServer(":8081")
	go httpServer.Run()

	grpcServer := server.NewGRPCServer(":8080")
	if err := grpcServer.Run(); err != nil {
		log.Fatal(err)
	}
}
