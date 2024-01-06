package main

import "google.golang.org/grpc"

const (
	grpcPort = "50051"
)

func main() {
	grpcServer := grpc.NewServer()
	orderService := UnimplementedOrderServiceServer{}
	RegisterOrderServiceServer(grpcServer, &orderService)

}
