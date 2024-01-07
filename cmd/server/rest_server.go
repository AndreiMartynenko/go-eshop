package main

import (
	"io"
	"net/http"

	"github.com/AndreiMartynenko/grpc-eshop/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// RestServer implements a REST server for the order service
type RestServer struct {
	server       *http.Server
	orderService proto.OrderServiceServer // The same order service as in the gRPC server
}

var router = gin.Default() // Declare a global router

// The NewRestServer function is perfect for creating a RestServer
func NewRestServer(orderService proto.OrderServiceServer, port string) RestServer {
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
	var req proto.CreateOrderRequest

	// Read the content of the request body into a []byte
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "error reading request body")
		return
	}

	// Request deserialization
	// err := protojson.Unmarshal(c.Request.Body, &req)
	err = protojson.Unmarshal(body, &req)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating order request")
		return
	}

	// Uses the order service to create an order from the request
	resp, err := r.orderService.Create(c.Request.Context(), &req)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating order")
		return
	}
	// Marshal the response to JSON
	//m := &protojsonpb.MarshalOptions{}
	m := &protojson.MarshalOptions{}
	if data, err := m.Marshal(resp); err != nil {
		c.String(http.StatusInternalServerError, "error sending order response")
	} else {
		c.Data(http.StatusOK, "application/json", data)
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
