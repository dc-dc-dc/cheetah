package market

import (
	"context"

	"github.com/dc-dc-dc/cheetah/util"
)

var _ MarketReceiver = (*ChainedReceiver)(nil)

const (
	ContextCache = "receiver.cache"
)

type ChainedReceiver struct {
	receivers []MarketReceiver
}

type MarketCache map[string]interface{}

func NewChainedReceiver(receivers ...MarketReceiver) *ChainedReceiver {
	cr := &ChainedReceiver{
		receivers: receivers,
	}
	cr.DedupReceivers(util.NewSet())
	return cr
}

func (cr *ChainedReceiver) DedupReceivers(keySet *util.Set) {
	res := []MarketReceiver{}
	for _, receiver := range cr.receivers {
		switch receiver.(type) {
		case CachableReceiver:
			{
				cachable := receiver.(CachableReceiver)
				if keySet.Contains(cachable.CacheKey()) {
					break
				}
				keySet.Add(cachable.CacheKey())
				res = append(res, receiver)
				break
			}
		case *ChainedReceiver:
			{
				chained := receiver.(*ChainedReceiver)
				chained.DedupReceivers(keySet)
				res = append(res, receiver)
			}
		default:
			res = append(res, receiver)
		}
	}
	cr.receivers = res
}

func (r *ChainedReceiver) Receive(ctx context.Context, line MarketLine) error {
	ctx = context.WithValue(ctx, ContextCache, make(MarketCache))
	for _, receiver := range r.receivers {
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}
