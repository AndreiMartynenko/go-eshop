package main

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
)

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

	// Route registration
	router.POST("/order", rs.create)
	router.GET("/order/:id", rs.retrieve)
	router.PUT("/order", rs.update)
	router.DELETE("/order", rs.delete)
	router.GET("/order", rs.list)

	return rs
}

// Start launches the server
func (r RestServer) Start() error {
	return r.server.ListenAndServe()
}

// The create handler function creates an order from the request (JSON body)
func (r RestServer) create(c *gin.Context) {
	var req CreateOrderRequest

	// Request deserialization
	err := jsonpb.Unmarshal(c.Request.Body, &req)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating order request")
	}

	// Uses the order service to create an order from the request
	resp, err := r.orderService.Create(c.Request.Context(), &req)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating order")
	}
	m := &jsonpb.Marshaler{}
	if err := m.Marshal(c.Writer, resp); err != nil {
		c.String(http.StatusInternalServerError, "error sending order response")
	}
}

func (r RestServer) retrieve(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented yet")
}

func (r RestServer) update(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented yet")
}

func (r RestServer) delete(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented yet")
}

func (r RestServer) list(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented yet")
}
