package market

import (
	"time"

	"github.com/shopspring/decimal"
)

type MarketLine struct {
	Start  time.Time
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume int64
}

// Note: This will panic if the string is not a valid decimal.
func NewMarketLineFromString(start time.Time, open, high, low, close string, vol int64) MarketLine {
	return MarketLine{
		Start:  start,
		Open:   decimal.RequireFromString(open),
		High:   decimal.RequireFromString(high),
		Low:    decimal.RequireFromString(low),
		Close:  decimal.RequireFromString(close),
		Volume: vol,
	}
}
