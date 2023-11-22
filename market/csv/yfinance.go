package csv

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
)

type YFinanceInterval int

const (
	YFinanceInterval1Minute YFinanceInterval = iota
	YFinanceInterval2Minute
	YFinanceInterval5Minute
	YFinanceInterval15Minute
	YFinanceInterval30Minute
	YFinanceInterval60Minute
	YFinanceInterval90Minute
	YFinanceInterval1Hour
	YFinanceInterval1Day
	YFinanceInterval5Day
	YFinanceInterval1Week
	YFinanceInterval1Month
	YFinanceInterval3Month
)

func (i YFinanceInterval) String() string {
	switch i {
	case YFinanceInterval1Minute:
		return "1m"
	case YFinanceInterval2Minute:
		return "2m"
	case YFinanceInterval5Minute:
		return "5m"
	case YFinanceInterval15Minute:
		return "15m"
	case YFinanceInterval30Minute:
		return "30m"
	case YFinanceInterval60Minute:
		return "60m"
	case YFinanceInterval90Minute:
		return "90m"
	case YFinanceInterval1Hour:
		return "1h"
	case YFinanceInterval1Day:
		return "1d"
	case YFinanceInterval5Day:
		return "5d"
	case YFinanceInterval1Week:
		return "1wk"
	case YFinanceInterval1Month:
		return "1mo"
	case YFinanceInterval3Month:
		return "3mo"
	}
	return ""
}

type yFinanceProducer struct {
	symbol   string
	interval YFinanceInterval
	start    time.Time
	end      time.Time
}

func NewYFinanceProducer(symbol string, interval YFinanceInterval, start, end time.Time) market.MarketProducer {
	return &yFinanceProducer{
		symbol:   symbol,
		interval: interval,
		start:    start,
		end:      end,
	}
}

func (p *yFinanceProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=%s&events=history", p.symbol, p.start.Unix(), p.end.Unix(), p.interval.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	defer res.Body.Close()
	return NewCsvProducer(res.Body).Produce(ctx, out)
}
