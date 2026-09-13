package main

func main() {
	httpServer := NewHttpServer(":8081")
	go httpServer.Run()

	grpcServer := NewGRPCServer(":8080")
	grpcServer.Run()

}
