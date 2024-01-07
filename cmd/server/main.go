package main

import (
	"log"
	"net"

	"github.com/AndreiMartynenko/grpc-eshop/proto"
	"google.golang.org/grpc"
)

const (
	grpcPort = "50051"
	restPort = "8080"
)

func main() {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	// Create an instance of the OrderServiceServer implementation
	orderService := proto.UnimplementedOrderServiceServer{}
	//orderService := &OrderServiceImpl{}
	// Register the OrderServiceServer with the gRPC server
	proto.RegisterOrderServiceServer(grpcServer, &orderService)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	go func() {
		// Serve()

	}()

}
