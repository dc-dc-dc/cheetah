package market

import (
	"context"
)

var _ MarketReceiver = (*chainedReceiver)(nil)

type chainedReceiver struct {
	receivers []MarketReceiver
}

func simplifyReceivers(receivers []MarketReceiver) []MarketReceiver {
	cache := make(map[string]interface{})
	res := []MarketReceiver{}
	for _, receiver := range receivers {
		if cachable, ok := receiver.(CachableReceiver); ok {
			if _, ok := cache[cachable.CacheKey()]; ok {
				continue
			}
			cache[cachable.CacheKey()] = nil
		}
		res = append(res, receiver)
	}
	return res
}

func NewChainedReceiver(receivers ...MarketReceiver) *chainedReceiver {
	return &chainedReceiver{
		receivers: simplifyReceivers(receivers),
	}
}

func (r *chainedReceiver) Receive(ctx context.Context, line MarketLine) error {
	ctx = context.WithValue(ctx, ContextCache, make(map[string]interface{}))
	for _, receiver := range r.receivers {
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}
