package main

import (
	"log"
	"net"

	"github.com/AndreiMartynenko/grpc-eshop/proto"
	"google.golang.org/grpc"
)

const (
	grpcPort = "50051"
)

func main() {
	grpcServer := grpc.NewServer()
	orderService := UnimplementedOrderServiceServer{}
	proto.RegisterOrderServiceServer(grpcServer, &orderService)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

}
