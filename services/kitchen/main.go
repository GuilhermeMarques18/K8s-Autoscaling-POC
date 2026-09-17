package main

import (
	"log"
	"os"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/kitchen/server"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	httpAddr := getEnv("HTTP_ADDR", ":9000")
	ordersGrpcAddr := getEnv("ORDERS_GRPC_ADDR", "localhost:8080")

	httpServer := server.NewHttpServer(httpAddr, ordersGrpcAddr)
	log.Fatal(httpServer.Run())
}
