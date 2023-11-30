package indicator

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/shopspring/decimal"
)

type macdReceiver market.MarketReceiver

func init() {
	market.RegisterSerializableReceiver(MacdCacheKey(), func() market.MarketReceiver {
		return newMacdReceiver()
	})
}

func GetMacdFromCache(ctx context.Context) (decimal.Decimal, error) {
	return market.GetFromCache[decimal.Decimal](ctx, MacdCacheKey())
}

func MacdCacheKey() string {
	return "indicator.macd"
}

func NewMacd() macdReceiver {
	return market.NewChainedReceiver(
		NewExponentialMovingAverage(12),
		NewExponentialMovingAverage(26),
		newMacdReceiver(),
	)
}

func newMacdReceiver() market.MarketReceiver {
	return market.NewCachableFunctionalReceiver(MacdCacheKey(), func(ctx context.Context, line market.MarketLine) error {
		ema12, err1 := GetExponentialMovingAverageFromCache(ctx, 12)
		ema26, err2 := GetExponentialMovingAverageFromCache(ctx, 26)
		if err1 != nil || err2 != nil {
			return nil
		}
		market.SetCache(ctx, MacdCacheKey(), ema12.Sub(ema26))
		return nil
	})
}
