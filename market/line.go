package market

import (
	"time"

	"github.com/shopspring/decimal"
)

type MarketLine struct {
	Start  time.Time       `json:"start"`
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume int64           `json:"volume"`
}

func EnsureMarketLineFromString(start time.Time, open, high, low, close string, vol int64) MarketLine {
	line, err := NewMarketLineFromString(start, open, high, low, close, vol)
	if err != nil {
		panic(err)
	}
	return *line
}

func IsMarketNull(line MarketLine) bool {
	return line.Open.IsZero() && line.High.IsZero() && line.Low.IsZero() && line.Close.IsZero() && line.Volume == 0
}

func NewMarketLineFromString(start time.Time, open, high, low, close string, vol int64) (*MarketLine, error) {
	var _open, _close, _high, _low decimal.Decimal
	var err error
	if _open, err = decimal.NewFromString(open); err != nil {
		return nil, err
	}
	if _close, err = decimal.NewFromString(close); err != nil {
		return nil, err
	}
	if _high, err = decimal.NewFromString(high); err != nil {
		return nil, err
	}
	if _low, err = decimal.NewFromString(low); err != nil {
		return nil, err
	}

	return &MarketLine{
		Start:  start,
		Open:   _open,
		High:   _high,
		Low:    _low,
		Close:  _close,
		Volume: vol,
	}, nil
}
