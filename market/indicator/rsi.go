package indicator

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

func init() {
	market.RegisterSerializableReceiver(RsiCacheKey(), func() market.MarketReceiver {
		return NewRsi()
	})
}

func RsiCacheKey() string {
	return "indicator.rsi"
}

func GetRsiFromCache(ctx context.Context) (decimal.Decimal, error) {
	return market.GetFromCache[decimal.Decimal](ctx, RsiCacheKey())
}

type rsiReceiver struct {
	last decimal.Decimal
	gain *util.CappedQueue[decimal.Decimal]
	loss *util.CappedQueue[decimal.Decimal]
}

func NewRsi() *rsiReceiver {
	return &rsiReceiver{
		gain: util.NewCappedQueue[decimal.Decimal](14),
		loss: util.NewCappedQueue[decimal.Decimal](14),
	}
}

func (r *rsiReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	if r.last.IsZero() {
		r.last = line.Close
		return nil
	}
	diff := line.Close.Sub(r.last)
	r.last = line.Close
	if diff.GreaterThan(decimal.Zero) {
		r.gain.Push(diff)
		r.loss.Push(decimal.Zero)
	} else {
		r.loss.Push(diff.Abs())
		r.gain.Push(decimal.Zero)
	}
	if r.loss.Full() {
		avg_loss := CalculateAverage(r.loss)
		avg_gain := CalculateAverage(r.gain)

		rsi := decimal.NewFromInt(100).Sub(
			decimal.NewFromInt(100).Div(
				decimal.NewFromInt(1).Add(
					avg_gain.Div(avg_loss),
				),
			),
		)
		market.SetCache(ctx, r.CacheKey(), rsi)
	}
	return nil
}

func (r *rsiReceiver) CacheKey() string {
	return RsiCacheKey()
}

func (r *rsiReceiver) PrefixKey() string {
	return RsiCacheKey()
}

var _ market.CachableReceiver = (*rsiReceiver)(nil)
