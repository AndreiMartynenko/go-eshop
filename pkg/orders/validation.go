package orders

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	ErrPreAuthorizationTimeout = errors.New("pre-authorization request timeout")
	ErrInventoryRequestTimeout = errors.New("check inventory request timeout")
	ErrItemOutOfStock          = errors.New("sorry, one or more items in your order is out of stock")
)

// preAuthorizePayment performs pre-authorization of the payment method
// and returns an error. nil is returned for successful pre-authorization
func preAuthorizePayment(ctx context.Context, payment *PaymentMethod, orderAmount float32) error {
	// Costly authorization logic is performed here - for this example, we use sleep mode :-)
	// and return nil to indicate successful authorization
	timer := time.NewTimer(3 * time.Second)

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ErrPreAuthorizationTimeout
	}
}

// checkInventory returns a boolean value and an error indicating
// whether all items are in stock. (true, nil) is returned if
// all items are in stock, and no errors occurred
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
