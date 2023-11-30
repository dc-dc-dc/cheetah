package indicator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/shopspring/decimal"
)

const (
	exponentialMovingAveragePrefixKey = "indicator.moving_average_exponential"
)

func ExponentialMovingAverageCacheKey(window int) string {
	return fmt.Sprintf("%s.%d", exponentialMovingAveragePrefixKey, window)
}

func GetExponentialMovingAverageFromCache(ctx context.Context, window int) (decimal.Decimal, error) {
	return market.GetFromCache[decimal.Decimal](ctx, ExponentialMovingAverageCacheKey(window))
}

type exponentialMovingAverage struct {
	window int
	last   decimal.Decimal
}

func NewExponentialMovingAverage(window int) *exponentialMovingAverage {
	return &exponentialMovingAverage{
		window: window,
	}
}

func (sa *exponentialMovingAverage) Receive(ctx context.Context, line market.MarketLine) error {
	if sa.last.IsZero() {
		sa.last = line.Close
		market.SetCache(ctx, sa.CacheKey(), sa.last)
		return nil
	}

	sa.last = (line.Close.Mul(decimal.NewFromFloat32(2.0 / float32(sa.window+1)))).Add(sa.last.Mul(decimal.NewFromFloat32(1.0 - (2.0 / float32(sa.window+1)))))
	market.SetCache(ctx, sa.CacheKey(), sa.last)
	return nil
}

func (sa *exponentialMovingAverage) CacheKey() string {
	return ExponentialMovingAverageCacheKey(sa.window)
}

func (sa *exponentialMovingAverage) PrefixKey() string {
	return exponentialMovingAveragePrefixKey
}

func (sa *exponentialMovingAverage) String() string {
	return fmt.Sprintf("ExponentialMovingAverage{window=%d}", sa.window)
}

func (sa *exponentialMovingAverage) MarshalJSON() ([]byte, error) {
	return json.Marshal(simpleMovingAverageJson{
		Window: sa.window,
	})
}

func (sa *exponentialMovingAverage) UnmarshalJSON(data []byte) error {
	var raw simpleMovingAverageJson
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	sa.window = raw.Window
	return nil
}

var _ market.CachableReceiver = (*exponentialMovingAverage)(nil)
var _ market.SerializableReceiver = (*exponentialMovingAverage)(nil)
