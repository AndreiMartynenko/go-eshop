package main

import "net/http"

// RestServer implements a REST server for the order service
type RestServer struct {
	server       *http.Server
	orderService OrderServiceServer // The same order service as in the gRPC server
}
