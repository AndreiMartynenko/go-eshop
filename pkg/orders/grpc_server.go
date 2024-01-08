package orders

import (
	"net"

	"github.com/AndreiMartynenko/grpc-eshop/proto"
	"google.golang.org/grpc"
)

// GrpcServer implements a gRPC server for the order service
type GrpcServer struct {
	server   *grpc.Server
	errCh    chan error
	listener net.Listener
}

// NewGrpcServer function is excellent for creating a GrpcServer
func NewGrpcServer(service proto.OrderServiceServer, port string) (GrpcServer, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return GrpcServer{}, err
	}
	server := grpc.NewServer()
	proto.RegisterOrderServiceServer(server, service)

	return GrpcServer{
		server:   server,
		listener: lis,
		errCh:    make(chan error),
	}, nil
}

// Start launches the gRPC server in the background, sending errors to the error channel
func (g GrpcServer) Start() {
	go func() {
		g.errCh <- g.server.Serve(g.listener)
	}()
}

// Stop stops the server
func (g GrpcServer) Stop() {
	g.server.GracefulStop()
}

// Error returns the server's error channel
func (g GrpcServer) Error() chan error {
	return g.errCh
}
