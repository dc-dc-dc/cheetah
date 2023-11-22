package basic

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketReceiver = (*basicReceiver)(nil)
var _ market.MarketReceiver = (*errorReceiver)(nil)

type errorReceiver struct {
	current   int
	explodeOn int
}

func NewErrorReceiver(explodesOn int) *errorReceiver {
	return &errorReceiver{
		explodeOn: explodesOn,
	}
}

func (r *errorReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	r.current++
	logLine(ctx, line)
	if r.current == r.explodeOn {
		return fmt.Errorf("receiver error")
	}
	return nil
}

type basicReceiver struct {
}

func NewBasicReceiver() *basicReceiver {
	return &basicReceiver{}
}

func (r *basicReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	logLine(ctx, line)
	return nil
}

func logLine(ctx context.Context, line market.MarketLine) {
	index, ok := ctx.Value("receiver").(int)
	if ok {
		fmt.Printf("[receiver-%d] %v\n", index, line)
	}
}
