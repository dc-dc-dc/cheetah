package indicator

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

var _ market.MarketReceiver = (*SimpleMovingAverage)(nil)

const ContextIndicatorSimpleMovingAverage = "indicator.simple_moving_average"

type SimpleMovingAverage struct {
	queue *util.CappedQueue
}

func NewSimpleMovingAverage(count int) *SimpleMovingAverage {
	return &SimpleMovingAverage{
		queue: util.NewCappedQueue(count),
	}
}

func (sa *SimpleMovingAverage) Receive(ctx context.Context, line market.MarketLine) error {
	cache, ok := ctx.Value(market.ContextCache).(map[string]interface{})
	if !ok {
		return market.ErrNoContextCache
	}

	sa.queue.Push(line.Close)
	if sa.queue.Full() {
		elements := sa.queue.Elements()
		var sum decimal.Decimal
		for _, element := range elements {
			sum = sum.Add(element.(decimal.Decimal))
		}
		average := sum.Div(decimal.NewFromInt(int64(sa.queue.Cap())))
		cache[ContextIndicatorSimpleMovingAverage] = average
	}
	return nil
}
