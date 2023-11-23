package market

import (
	"context"
)

var _ MarketReceiver = (*chainedReceiver)(nil)

type chainedReceiver struct {
	receivers []MarketReceiver
}

func NewChainedReceiver(receivers ...MarketReceiver) *chainedReceiver {
	return &chainedReceiver{
		receivers: receivers,
	}
}

func (r *chainedReceiver) Receive(ctx context.Context, line MarketLine) error {
	for _, receiver := range r.receivers {
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}
