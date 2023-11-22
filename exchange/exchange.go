package exchange

import (
	"context"
	"errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Exchange interface {
	GetPositions(context.Context) ([]Position, error)
	PlaceOrder(context.Context, Order) error
	CancelOrder(context.Context, ID) error
}
