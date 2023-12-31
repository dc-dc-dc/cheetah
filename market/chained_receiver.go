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
	RegisterSerializableReceiver(chainedReceiverPrefixKey, func() MarketReceiver {
		return &ChainedReceiver{}
	})
}

type ChainedReceiver struct {
	receivers []MarketReceiver
}

func NewChainedReceiver(receivers ...MarketReceiver) *ChainedReceiver {
	cr := &ChainedReceiver{}
	cr.receivers = DedupReceivers(receivers, util.NewSet[string]())
	return cr
}

func (cr *ChainedReceiver) PrefixKey() string {
	return chainedReceiverPrefixKey
}

func (cr *ChainedReceiver) Receivers() []MarketReceiver {
	return cr.receivers
}

func DedupReceivers(receivers []MarketReceiver, keySet *util.Set[string]) []MarketReceiver {
	res := make([]MarketReceiver, 0, len(receivers))
	for _, receiver := range receivers {
		switch receiver := receiver.(type) {
		case CachableReceiver:
			{
				if keySet.Contains(receiver.CacheKey()) {
					break
				}
				keySet.Add(receiver.CacheKey())
				res = append(res, receiver)
				break
			}
		case *ChainedReceiver:
			{
				res = append(res, DedupReceivers(receiver.Receivers(), keySet)...)
			}
		default:
			res = append(res, receiver)
		}
	}
	return res
}

func (r *ChainedReceiver) Receive(ctx context.Context, line MarketLine) error {
	if _, err := GetCache(ctx); err != nil {
		// Only create the cache if it does not exist
		ctx = CreateCache(ctx)
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
