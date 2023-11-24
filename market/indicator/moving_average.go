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
	movingAveragePrefixKey = "indicator.moving_average"
)

func init() {
	market.RegisterSerializableReceiver(movingAveragePrefixKey, func() market.SerializableReceiver {
		return &MovingAverage{}
	})
}

type MovingAverageCalc func([]decimal.Decimal) decimal.Decimal

func SimpleMovingAverageCalc(items []decimal.Decimal) decimal.Decimal {
	var sum decimal.Decimal
	for _, element := range items {
		sum = sum.Add(element)
	}
	return sum.Div(decimal.NewFromInt(int64(len(items))))
}

func ExponentialMovingAverageCalc(items []decimal.Decimal) decimal.Decimal {
	var sum decimal.Decimal
	k := decimal.NewFromFloat32(2.0 / float32(len(items)+1))
	for i, element := range items {
		if i == 0 {
			sum = element
		} else {
			sum = element.Mul(k).Add(sum.Mul(decimal.NewFromFloat32(1.0).Sub(k)))
		}
	}
	return sum
}

func SimpleMovingAverageCacheKey(window int) string {
	return fmt.Sprintf("%s.%d.simple", movingAveragePrefixKey, window)
}

func ExponentialMovingAverageCacheKey(window int) string {
	return fmt.Sprintf("%s.%d.exponential", movingAveragePrefixKey, window)
}

type MovingAverage struct {
	queue  *util.CappedQueue[decimal.Decimal]
	simple bool
}

func NewSimpleMovingAverage(count int) *MovingAverage {
	return NewMovingAverage(count, true)
}

func NewExponentialMovingAverage(count int) *MovingAverage {
	return NewMovingAverage(count, false)
}

func NewMovingAverage(count int, simple bool) *MovingAverage {
	return &MovingAverage{
		queue:  util.NewCappedQueue[decimal.Decimal](count),
		simple: simple,
	}
}

func (sa *MovingAverage) CacheKey() string {
	if sa.simple {
		return SimpleMovingAverageCacheKey(sa.queue.Cap())
	}
	return ExponentialMovingAverageCacheKey(sa.queue.Cap())
}

func (sa *MovingAverage) PrefixKey() string {
	return movingAveragePrefixKey
}

func (sa *MovingAverage) Receive(ctx context.Context, line market.MarketLine) error {
	cache, err := market.GetCache(ctx)
	if err != nil {
		return err
	}
	sa.queue.Push(line.Close)
	if sa.queue.Full() {
		items := sa.queue.Elements()
		if sa.simple {
			cache[sa.CacheKey()] = SimpleMovingAverageCalc(items)
		} else {
			cache[sa.CacheKey()] = ExponentialMovingAverageCalc(items)
		}
	}
	return nil
}

func (sa *MovingAverage) String() string {
	return fmt.Sprintf("MovingAverage{window=%d, simple=%t}", sa.queue.Cap(), sa.simple)
}

type movingAverageJson struct {
	Simple bool `json:"simple"`
	Window int  `json:"window"`
}

func (sa *MovingAverage) MarshalJSON() ([]byte, error) {
	return json.Marshal(movingAverageJson{
		Simple: sa.simple,
		Window: sa.queue.Cap(),
	})
}

func (sa *MovingAverage) UnmarshalJSON(data []byte) error {
	var raw movingAverageJson
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	sa.simple = raw.Simple
	sa.queue = util.NewCappedQueue[decimal.Decimal](raw.Window)
	return nil
}

var _ market.CachableReceiver = (*MovingAverage)(nil)
var _ market.SerializableReceiver = (*MovingAverage)(nil)
var _ MovingAverageCalc = SimpleMovingAverageCalc
