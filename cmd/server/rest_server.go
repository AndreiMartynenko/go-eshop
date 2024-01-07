package main

import "net/http"

// RestServer implements a REST server for the order service
type RestServer struct {
	server       *http.Server
	orderService OrderServiceServer // The same order service as in the gRPC server
}


// The NewRestServer function is perfect for creating a RestServer
func NewRestServer(orderService OrderServiceServer, port string) RestServer {
	rs := RestServer{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		orderService: orderService,
	}