package market

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/util"
)

var _ SerializableReceiver = (*ChainedReceiver)(nil)

const (
	chainedReceiverPrefixKey = "chained_receiver"
)

func init() {
	RegisterSerializableReceiver(chainedReceiverPrefixKey, func() SerializableReceiver {
		return &ChainedReceiver{}
	})
}

type ChainedReceiver struct {
	receivers []MarketReceiver
}

func NewChainedReceiver(receivers ...MarketReceiver) *ChainedReceiver {
	cr := &ChainedReceiver{
		receivers: receivers,
	}
	cr.DedupReceivers(util.NewSet())
	return cr
}

func (cr *ChainedReceiver) PrefixKey() string {
	return chainedReceiverPrefixKey
}

func (cr *ChainedReceiver) Receivers() []MarketReceiver {
	return cr.receivers
}

func (cr *ChainedReceiver) DedupReceivers(keySet *util.Set) {
	res := make([]MarketReceiver, 0, len(cr.receivers))
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
	if _, err := GetCache(ctx); err != nil {
		// Only create the cache if it does not exist
		ctx = context.WithValue(ctx, ContextCache, make(MarketCache))
	}
	for _, receiver := range r.receivers {
		if err := receiver.Receive(ctx, line); err != nil {
			return err
		}
	}
	return nil
}

func (r *ChainedReceiver) MarshalJSON() ([]byte, error) {
	return SerializeReceivers(r.receivers...)
}

func (r *ChainedReceiver) UnmarshalJSON(data []byte) error {
	rec, err := DeserializeReceivers(data)
	if err != nil {
		return err
	}
	r.receivers = rec
	return nil
}

func (r *ChainedReceiver) String() string {
	return fmt.Sprintf("ChainedReceiver{receivers=%v}", r.receivers)
}
