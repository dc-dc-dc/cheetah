package indicator

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

var _ market.MarketReceiver = (*MovingAverage)(nil)
var _ MovingAverageCalc = SimpleMovingAverageCalc

const ContextIndicatorSimpleMovingAverage = "indicator.moving_average"

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
	queue *util.CappedQueue
	calc  MovingAverageCalc
}

func NewSimpleMovingAverage(count int) *MovingAverage {
	return NewMovingAverage(count, SimpleMovingAverageCalc)
}

func NewExponentialMovingAverage(count int) *MovingAverage {
	return NewMovingAverage(count, ExponentialMovingAverageCalc)
}

func NewMovingAverage(count int, calc MovingAverageCalc) *MovingAverage {
	return &MovingAverage{
		queue: util.NewCappedQueue(count),
		calc:  calc,
	}
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
		cache[ContextIndicatorSimpleMovingAverage] = sa.calc(items)
	}
	return nil
}
