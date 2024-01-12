package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

// OrderDispatcher is a daemon process that creates a set of handlers using sync.WaitGroup to concurrently
// process and dispatch orders
type OrderDispatcher struct {
	ordersCh   chan *orders.Order
	orderLimit int // maximum number of orders the pool will process concurrently
}

// NewOrderDispatcher creates a new OrderDispatcher
func NewOrderDispatcher(orderLimit int, bufferSize int) OrderDispatcher {
	return OrderDispatcher{
		ordersCh:   make(chan *orders.Order, bufferSize),
		orderLimit: orderLimit,
	}
}

// SubmitOrder submits an order for processing
func (d OrderDispatcher) SubmitOrder(order *orders.Order) {
	go func() {
		d.ordersCh <- order
	}()
}

// Start launches the dispatcher in the background
func (d OrderDispatcher) Start() {
	go d.processOrders()
}

// Shutdown shuts down the OrderDispatcher by closing the orders channel
// Note: This function should only be executed after the last order
// has entered the orders channel. Sending an order to a closed channel will panic.
func (d OrderDispatcher) Shutdown() {
	close(d.ordersCh)
}

// processOrders processes all incoming orders in the background using
// for-range and sync.WaitGroup
func (d OrderDispatcher) processOrders() {
	limiter := make(chan struct{}, d.orderLimit)
	var wg sync.WaitGroup

	// Continuous processing of orders received from the orders channel
	// This loop will exit after the channel is closed
	for order := range d.ordersCh {
		limiter <- struct{}{}
		wg.Add(1)

		go func(order *orders.Order) {
			// What needs to be done: start the fulfillment process to pack and ship the order
			// Currently using a sleep and print for demonstration
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Order (%v) has shipped \n", order)
			<-limiter
			wg.Done()
		}(order)
	}
	wg.Wait()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	dispatcher := NewOrderDispatcher(3, 100)
	dispatcher.Start()
	defer dispatcher.Shutdown()

	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "iPhone Screen Protector", Price: 9.99}}})
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "iPhone Case", Price: 19.99}}})
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "Pixel Case", Price: 14.99}}})
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "Bluetooth Speaker", Price: 29.99}}})
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "4K Monitor", Price: 159.99}}})
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "Inkjet Printer", Price: 79.99}}})

	time.Sleep(5 * time.Second) // just for testing

	// Continue submitting orders or perform any other operations
	// before the main function exits.
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "Mouse", Price: 14.99}}})
	dispatcher.SubmitOrder(&orders.Order{Items: []*orders.Item{{Description: "Keyboard", Price: 29.99}}})

	// Let the program run for a while to process the additional orders
	time.Sleep(5 * time.Second)

	// After processing, you can shut down the dispatcher
	dispatcher.Shutdown()
}
