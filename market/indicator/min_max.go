package indicator

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

var _ market.CachableReceiver = (*MinMaxIndicator)(nil)

type indexPrice struct {
	index int
	price decimal.Decimal
}

type MinMaxIndicator struct {
	count  int
	queue  *util.Queue
	window int
	min    bool
}

func NewMinIndicator(window int) *MinMaxIndicator {
	return newMinMaxIndicator(window, true)
}

func NewMaxIndicator(window int) *MinMaxIndicator {
	return newMinMaxIndicator(window, false)
}

func newMinMaxIndicator(window int, min bool) *MinMaxIndicator {
	return &MinMaxIndicator{
		queue:  util.NewQueue(),
		min:    min,
		window: window,
	}
}

func (mm *MinMaxIndicator) CacheKey() string {
	return fmt.Sprintf("indicator.min_max.%d.%s", mm.window, mm.Type())
}

func (mm *MinMaxIndicator) Type() string {
	if mm.min {
		return "min"
	}
	return "max"
}

func (mm *MinMaxIndicator) compare(first, other decimal.Decimal) bool {
	if mm.min {
		return first.GreaterThan(other)
	}
	return first.LessThan(other)
}
func (mm *MinMaxIndicator) Receive(ctx context.Context, line market.MarketLine) error {
	mm.count += 1
	cache, ok := ctx.Value(market.ContextCache).(market.MarketCache)
	if !ok {
		return market.ErrNoContextCache
	}
	for mm.queue.Count() > 0 && mm.count >= mm.queue.First().(indexPrice).index {
		mm.queue.PopLeft()
	}

	for mm.queue.Count() > 0 && mm.compare(mm.queue.Last().(indexPrice).price, line.Close) {
		mm.queue.Pop()
	}
	mm.queue.Push(indexPrice{mm.count + mm.window, line.Close})
	cache[mm.CacheKey()] = mm.queue.First().(indexPrice).price
	return nil
}
