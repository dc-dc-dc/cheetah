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
	minMaxPrefixKey = "indicator.min_max"
)

func init() {
	market.RegisterSerializableReceiver(minMaxPrefixKey, func() market.SerializableReceiver {
		return &MinMaxIndicator{}
	})
}

var _ market.CachableReceiver = (*MinMaxIndicator)(nil)
var _ market.SerializableReceiver = (*MinMaxIndicator)(nil)

func MinMaxCacheKey(window int, min bool) string {
	if min {
		return fmt.Sprintf("%s.%d.min", minMaxPrefixKey, window)
	}
	return fmt.Sprintf("%s.%d.max", minMaxPrefixKey, window)
}

type indexPrice struct {
	index int
	price decimal.Decimal
}

type MinMaxIndicator struct {
	count  int
	queue  *util.Queue[indexPrice]
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
		queue:  util.NewQueue[indexPrice](),
		min:    min,
		window: window,
	}
}

func (mm *MinMaxIndicator) PrefixKey() string {
	return minMaxPrefixKey
}

func (mm *MinMaxIndicator) CacheKey() string {
	return MinMaxCacheKey(mm.window, mm.min)
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
	for mm.queue.Count() > 0 && mm.count >= mm.queue.First().index {
		mm.queue.PopLeft()
	}

	for mm.queue.Count() > 0 && mm.compare(mm.queue.Last().price, line.Close) {
		mm.queue.Pop()
	}
	mm.queue.Push(indexPrice{mm.count + mm.window, line.Close})
	cache[mm.CacheKey()] = mm.queue.First().price
	return nil
}

type minMaxIndicatorJSON struct {
	Window int  `json:"window"`
	Min    bool `json:"min"`
}

func (mm *MinMaxIndicator) UnmarshalJSON(data []byte) error {
	var j minMaxIndicatorJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	mm.window = j.Window
	mm.min = j.Min
	mm.queue = util.NewQueue[indexPrice]()
	return nil
}

func (mm *MinMaxIndicator) MarshalJSON() ([]byte, error) {
	return json.Marshal(minMaxIndicatorJSON{
		Window: mm.window,
		Min:    mm.min,
	})
}

func (mm *MinMaxIndicator) String() string {
	return fmt.Sprintf("MinMaxIndicator{window=%d, min=%v}", mm.window, mm.min)
}
