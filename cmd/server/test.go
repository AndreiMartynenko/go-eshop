package main

import "github.com/AndreiMartynenko/grpc-eshop/proto"

// Avoid this
type OrderTotaler struct {
	items []*proto.Item
}

// This is a method. Binding it to a structure brings no benefit
// The structure must be initialized before testing
// this method
func (t OrderTotaler) getOrderTotal() float32 {
	var total float32

	for _, item := range t.items {
		total += item.Price
	}

	return total
}

// Prefer this approach. It's a pure function
func getOrderTotal(items []*proto.Item) float32 {
	var total float32

	for _, item := range items {
		total += item.Price
	}

	return total
}
