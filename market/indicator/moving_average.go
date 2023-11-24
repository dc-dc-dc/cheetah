package indicator

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

var _ market.CachableReceiver = (*MovingAverage)(nil)
var _ MovingAverageCalc = SimpleMovingAverageCalc

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

type MovingAverage struct {
	queue  *util.CappedQueue
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
		queue:  util.NewCappedQueue(count),
		simple: simple,
	}
}

func (sa *MovingAverage) CacheKey() string {
	return fmt.Sprintf("indicator.moving_average.%s.%d", sa.Type(), sa.queue.Cap())
}

func (sa *MovingAverage) Type() string {
	if sa.simple {
		return "simple"
	}
	return "exponential"
}

func (sa *MovingAverage) Receive(ctx context.Context, line market.MarketLine) error {
	cache, ok := ctx.Value(market.ContextCache).(map[string]interface{})
	if !ok {
		return market.ErrNoContextCache
	}

	sa.queue.Push(line.Close)
	if sa.queue.Full() {
		items := make([]decimal.Decimal, sa.queue.Cap())
		for i, t := range sa.queue.Elements() {
			items[i] = t.(decimal.Decimal)
		}
		if sa.simple {
			cache[sa.CacheKey()] = SimpleMovingAverageCalc(items)
		} else {
			cache[sa.CacheKey()] = ExponentialMovingAverageCalc(items)
		}
	}
	return nil
}
