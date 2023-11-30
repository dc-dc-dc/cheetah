package indicator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

const (
	simpleMovingAveragePrefixKey = "indicator.moving_average_simple"
)

func init() {
	market.RegisterSerializableReceiver(simpleMovingAveragePrefixKey, func() market.MarketReceiver {
		return &simpleMovingAverage{}
	})
}

func SimpleMovingAverageCacheKey(window int) string {
	return fmt.Sprintf("%s.%d", simpleMovingAveragePrefixKey, window)
}

type simpleMovingAverage struct {
	queue *util.CappedQueue[decimal.Decimal]
}

func NewSimpleMovingAverage(count int) *simpleMovingAverage {
	return &simpleMovingAverage{
		queue: util.NewCappedQueue[decimal.Decimal](count),
	}
}

func (sa *simpleMovingAverage) CacheKey() string {
	return SimpleMovingAverageCacheKey(sa.queue.Cap())
}

func (sa *simpleMovingAverage) PrefixKey() string {
	return simpleMovingAveragePrefixKey
}

func (sa *simpleMovingAverage) Receive(ctx context.Context, line market.MarketLine) error {
	sa.queue.Push(line.Close)
	items := sa.queue.Elements()
	var sum decimal.Decimal
	for _, element := range items {
		sum = sum.Add(element)
	}
	market.SetCache(ctx, sa.CacheKey(), sum.Div(decimal.NewFromInt(int64(len(items)))))
	return nil
}

func (sa *simpleMovingAverage) String() string {
	return fmt.Sprintf("SimpleMovingAverage{window=%d}", sa.queue.Cap())
}

type simpleMovingAverageJson struct {
	Window int `json:"window"`
}

func (sa *simpleMovingAverage) MarshalJSON() ([]byte, error) {
	return json.Marshal(simpleMovingAverageJson{
		Window: sa.queue.Cap(),
	})
}

func (sa *simpleMovingAverage) UnmarshalJSON(data []byte) error {
	var raw simpleMovingAverageJson
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	sa.queue = util.NewCappedQueue[decimal.Decimal](raw.Window)
	return nil
}

var _ market.CachableReceiver = (*simpleMovingAverage)(nil)
var _ market.SerializableReceiver = (*simpleMovingAverage)(nil)
