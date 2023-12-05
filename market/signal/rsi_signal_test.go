package signal_test

import (
	"context"
	"testing"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/indicator"
	"github.com/dc-dc-dc/cheetah/market/signal"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRsiSignal(t *testing.T) {
	rsiSignal := signal.NewRsiSignal()
	assert.Equal(t, "signal.rsi", rsiSignal.CacheKey())
	assert.Equal(t, "signal.rsi", rsiSignal.PrefixKey())
	ctx := market.CreateCache(context.Background())
	assert.NoError(t, rsiSignal.Receive(ctx, market.MarketLine{}))
	_, err := signal.GetRsiSignalFromCache(ctx)
	assert.Error(t, err)

	market.SetCache(ctx, indicator.RsiCacheKey(), decimal.NewFromFloat(80))
	assert.NoError(t, rsiSignal.Receive(ctx, market.MarketLine{}))
	market.SetCache(ctx, indicator.RsiCacheKey(), decimal.NewFromFloat(20))
	assert.NoError(t, rsiSignal.Receive(ctx, market.MarketLine{}))
	v1, err := signal.GetRsiSignalFromCache(ctx)
	assert.NoError(t, err)
	assert.Equal(t, signal.BuySignal, v1)
	ctx = market.CreateCache(context.Background())
	market.SetCache(ctx, indicator.RsiCacheKey(), decimal.NewFromFloat(20))

	assert.NoError(t, rsiSignal.Receive(ctx, market.MarketLine{}))
	_, err = signal.GetRsiSignalFromCache(ctx)
	assert.Error(t, err)

	market.SetCache(ctx, indicator.RsiCacheKey(), decimal.NewFromFloat(80))
	assert.NoError(t, rsiSignal.Receive(ctx, market.MarketLine{}))
	v2, err := signal.GetRsiSignalFromCache(ctx)
	assert.NoError(t, err)
	assert.Equal(t, signal.SellSignal, v2)
}
