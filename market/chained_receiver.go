package market

import (
	"context"
)

var _ MarketReceiver = (*chainedReceiver)(nil)

const (
	ContextCache = "receiver.cache"
)

type chainedReceiver struct {
	receivers []MarketReceiver
}

type MarketCache map[string]interface{}

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
	ctx = context.WithValue(ctx, ContextCache, make(MarketCache))
	for _, receiver := range r.receivers {
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}
