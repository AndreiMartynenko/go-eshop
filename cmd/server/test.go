package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/AndreiMartynenko/grpc-eshop/pkg/orders"
	"github.com/AndreiMartynenko/grpc-eshop/proto"
)

// MockOrderServiceServer implements the OrderServiceServer interface for testing
type MockOrderServiceServer struct{}

func (m MockOrderServiceServer) Create() {
	// Implement Create method for testing
}

// MockRestServer implements the RestServer interface for testing
type MockRestServer struct {
	startCalled bool
	stopCalled  bool
}

func (m *MockRestServer) Start() {
	m.startCalled = true
}

func (m *MockRestServer) Stop() error {
	m.stopCalled = true
	return nil
}

func (m MockRestServer) Error() chan error {
	return make(chan error)
}

// MockGrpcServer implements the GrpcServer interface for testing
type MockGrpcServer struct {
	startCalled bool
	stopCalled  bool
}

func (m *MockGrpcServer) Start() {
	m.startCalled = true
}

func (m *MockGrpcServer) Stop() {
	m.stopCalled = true
}

func (m MockGrpcServer) Error() chan error {
	return make(chan error)
}

func TestAppStartAndShutdown(t *testing.T) {
	// Create an app with mock servers
	mockOrderService := MockOrderServiceServer{}
	mockRestServer := &MockRestServer{}
	mockGrpcServer := &MockGrpcServer{}
	app := app{
		restServer: mockRestServer,
		grpcServer: mockGrpcServer,
		shutdownCh: make(chan os.Signal, 1),
	}

	// Start the app
	app.start()

	// Allow some time for servers to start (non-blocking)
	time.Sleep(100 * time.Millisecond)

	// Check if start was called on both servers
	if !mockRestServer.startCalled {
		t.Error("RestServer Start method not called")
	}
	if !mockGrpcServer.startCalled {
		t.Error("GrpcServer Start method not called")
	}

	// Stop the app
	err := app.shutdown()

	// Check if stop was called on both servers
	if !mockRestServer.stopCalled {
		t.Error("RestServer Stop method not called")
	}
	if !mockGrpcServer.stopCalled {
		t.Error("GrpcServer Stop method not called")
	}

	// Check if the shutdown error is nil
	if err != nil {
		t.Errorf("Expected shutdown error to be nil, got %v", err)
	}
}

func TestRunFunctionWithError(t *testing.T) {
	// Override newApp to return an error
	newApp = func() (app, error) {
		return app{}, errors.New("error creating app")
	}

	// Run the application
	err := run()

	// Check if the run function returns an error
	if err == nil {
		t.Error("Expected run function to return an error, got nil")
	}

	// Reset newApp to the original implementation
	newApp = func() (app, error) {
		orderService := proto.UnimplementedOrderServiceServer{}
		gs, _ := orders.NewGrpcServer(orderService, grpcPort)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		return app{
			restServer: orders.NewRestServer(orderService, restPort),
			grpcServer: gs,
			shutdownCh: quit,
		}, nil
	}
}

// Add more test cases as needed for specific components
