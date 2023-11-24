package indicator

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/shopspring/decimal"
)

type macdReceiver = market.MarketReceiver

func MacdCacheKey() string {
	return "indicator.macd"
}

func NewMacd() macdReceiver {
	return market.NewChainedReceiver(
		NewExponentialMovingAverage(12),
		NewExponentialMovingAverage(26),
		market.NewFunctionalReceiver(func(ctx context.Context, line market.MarketLine) error {
			cache, err := market.GetCache(ctx)
			if err != nil {
				return market.ErrNoContextCache
			}
			// TODO: Think of a better way of retreiving from the cache
			ema12, ok1 := cache[ExponentialMovingAverageCacheKey(12)]
			ema26, ok2 := cache[ExponentialMovingAverageCacheKey(26)]
			if !ok1 || !ok2 {
				return nil
			}
			ema12Val, ok1 := ema12.(decimal.Decimal)
			ema26Val, ok2 := ema26.(decimal.Decimal)
			if !ok1 || !ok2 {
				return fmt.Errorf("invalid ema12 or ema26")
			}
			cache[MacdCacheKey()] = ema12Val.Sub(ema26Val)
			return nil
		}),
	)
}
