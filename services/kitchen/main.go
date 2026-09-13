package main

import (
	"log"

	"github.com/GuilhermeMarques18/K8s-Autoscaling-POC.git/services/kitchen/server"
)

func main() {
	httpServer := server.NewHttpServer(":9000", ":8080")
	log.Fatal(httpServer.Run())
}
