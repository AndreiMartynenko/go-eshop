package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreiMartynenko/grpc-eshop/pkg/orders"
	"github.com/AndreiMartynenko/grpc-eshop/proto"
)

const (
	grpcPort = "50051"
	restPort = "8080"
)

// The app wrapper is perfect for all elements needed to start
// and stop the Order microservice
type app struct {
	restServer orders.RestServer
	grpcServer orders.GrpcServer
	//Listens for an application termination signal
	//Ex. (Ctrl X, Docker container shutdown, etc)
	shutdownCh chan os.Signal
}

// start launches the REST and gRPC servers in the background
func (a app) start() {
	go a.restServer.Start() // non-blocking now
	go a.grpcServer.Start() // also non-blocking :-)
}

// stop stops the servers
func (a app) shutdown() error {
	a.grpcServer.Stop()
	return a.restServer.Stop()
}

// newApp creates a new application with REST and gRPC servers
// This function performs all necessary application initialization
func newApp() (app, error) {
	orderService := proto.UnimplementedOrderServiceServer{}

	gs, err := orders.NewGrpcServer(orderService, grpcPort)
	if err != nil {
		return app{}, err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	return app{
		restServer: orders.NewRestServer(orderService, restPort),
		grpcServer: gs,
		shutdownCh: quit,
	}, nil
}

// run starts the application, handling any errors from REST and gRPC servers
// and shutdown signals
func run() error {
	app, err := newApp()
	if err != nil {
		return err
	}

	app.start()
	defer app.shutdown()

	select {
	case restErr := <-app.restServer.Error():
		return restErr
	case grpcErr := <-app.grpcServer.Error():
		return grpcErr
	case <-app.shutdownCh:
		return nil
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
