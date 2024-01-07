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
	// Listen for gRPC requests on the specified port
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
	// Serve gRPC requests in a separate goroutine
	go func() {
		// Serve()
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}

	}()

}
