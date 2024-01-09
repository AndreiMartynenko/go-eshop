package orders

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	ErrPreAuthorizationTimeout = errors.New("pre-authorization request timeout")
	ErrInventoryRequestTimeout = errors.New("check inventory request timeout")
	ErrItemOutOfStock          = errors.New("sorry, one or more items in your order is out of stock")
)

func preAuthorizePayment(ctx context.Context, payment *PaymentMethod, orderAmount float32) error {
	// Your pre-authorization logic
}

func checkInventory(ctx context.Context, items []*Item) (bool, error) {
	// Your inventory checking logic
}

func getOrderTotal(items []*Item) float32 {
	// Your logic to calculate order total
}

func validateOrder(ctx context.Context, items []*Item, payment *PaymentMethod) error {
	g, errCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return preAuthorizePayment(errCtx, payment, getOrderTotal(items))
	})

	g.Go(func() error {
		itemsInStock, err := checkInventory(errCtx, items)
		if err != nil {
			return err
		}
		if !itemsInStock {
			return ErrItemOutOfStock
		}
		return nil
	})

	return g.Wait()
}
