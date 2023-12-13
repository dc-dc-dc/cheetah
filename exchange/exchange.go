package exchange

import (
	"context"
	"errors"

	"github.com/dc-dc-dc/cheetah/util"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Exchange interface {
	GetPositions(context.Context) ([]Position, error)
	PlaceOrder(context.Context, Order) error
	CancelOrder(context.Context, util.ID) error
}
