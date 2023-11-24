package indicator

import (
	"context"

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
			ema12, err1 := market.GetFromCache[decimal.Decimal](ctx, ExponentialMovingAverageCacheKey(12))
			ema26, err2 := market.GetFromCache[decimal.Decimal](ctx, ExponentialMovingAverageCacheKey(26))
			if err1 != nil || err2 != nil {
				return nil
			}
			market.SetCache(ctx, MacdCacheKey(), ema12.Sub(ema26))
			return nil
		}),
	)
}
