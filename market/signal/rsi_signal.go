package signal

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/indicator"
	"github.com/shopspring/decimal"
)

type RsiSignal struct {
	last decimal.Decimal
}

func NewRsiSignal() *RsiSignal {
	return &RsiSignal{}
}

func (r *RsiSignal) Receive(ctx context.Context, line market.MarketLine) error {
	t, err := indicator.GetRsiFromCache(ctx)
	if err != nil {
		return nil
	}
	if r.last.IsZero() {
		r.last = t
		return nil
	}
	if t.GreaterThan(decimal.NewFromInt(70)) && r.last.LessThan(decimal.NewFromInt(70)) {
		market.SetCache(ctx, r.CacheKey(), SellSignal)
	}
	if t.LessThan(decimal.NewFromInt(30)) && r.last.GreaterThan(decimal.NewFromInt(30)) {
		market.SetCache(ctx, r.CacheKey(), BuySignal)
	}
	r.last = t
	return nil
}

func (r *RsiSignal) CacheKey() string {
	return "signal.rsi"
}

func (r *RsiSignal) PrefixKey() string {
	return "signal.rsi"
}

var _ market.CachableReceiver = (*RsiSignal)(nil)
