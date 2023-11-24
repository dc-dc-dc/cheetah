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
	cache, ok := ctx.Value(ContextCache).(map[string]interface{})
	for _, receiver := range r.receivers {
		if ok {
			// check if the receiver is cachable and has been cached
			if cr, ok := receiver.(CachableReceiver); ok {
				if _, ok := cache[cr.CacheKey()]; ok {
					continue
				}
			}
		}
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}
