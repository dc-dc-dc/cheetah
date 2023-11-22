package basic

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketReceiver = (*chainedReceiver)(nil)

type chainedReceiver struct {
	receivers []market.MarketReceiver
}

func NewChainedReceiver(receivers ...market.MarketReceiver) *chainedReceiver {
	return &chainedReceiver{
		receivers: receivers,
	}
}

func (r *chainedReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	for _, receiver := range r.receivers {
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}
