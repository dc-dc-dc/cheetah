package indicator

import (
	"context"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/shopspring/decimal"
)

func init() {
	market.RegisterSerializableReceiver(MacdCacheKey(), func() market.MarketReceiver {
		return newMacdReceiver()
	})
}

var (
	macdSignalConst = decimal.NewFromFloat32(2.0 / 10.0)
)

type MacdData struct {
	Macd      decimal.Decimal `json:"macd"`
	Signal    decimal.Decimal `json:"signal"`
	Histogram decimal.Decimal `json:"histogram"`
}

func GetMacdFromCache(ctx context.Context) (MacdData, error) {
	return market.GetFromCache[MacdData](ctx, MacdCacheKey())
}

func MacdCacheKey() string {
	return "indicator.macd"
}

func NewMacd() market.MarketReceiver {
	return market.NewChainedReceiver(
		NewExponentialMovingAverage(12),
		NewExponentialMovingAverage(26),
		newMacdReceiver(),
	)
}

type macdReceiver struct {
	last decimal.Decimal
}

func newMacdReceiver() *macdReceiver {
	return &macdReceiver{}
}

func (m *macdReceiver) CacheKey() string {
	return MacdCacheKey()
}

func (m *macdReceiver) PrefixKey() string {
	return MacdCacheKey()
}

func (m *macdReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	ema12, err1 := GetExponentialMovingAverageFromCache(ctx, 12)
	ema26, err2 := GetExponentialMovingAverageFromCache(ctx, 26)
	if err1 != nil || err2 != nil {
		return nil
	}
	macd := ema12.Sub(ema26)
	signal := (macd.Mul(macdSignalConst)).Add(m.last.Mul(decimal.NewFromFloat32(1.0).Sub(macdSignalConst)))
	m.last = signal

	market.SetCache(ctx, MacdCacheKey(), MacdData{
		Macd:      macd,
		Signal:    signal,
		Histogram: macd.Sub(signal),
	})
	return nil
}

var _ market.CachableReceiver = (*macdReceiver)(nil)
