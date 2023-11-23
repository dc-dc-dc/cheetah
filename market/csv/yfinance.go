package csv

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
)

type yFinanceProducer struct {
	symbol   market.Symbol
	interval market.Interval
	start    time.Time
	end      time.Time
}

func NewYFinanceProducer(symbol string, interval market.Interval, start, end time.Time) market.MarketProducer {
	return &yFinanceProducer{
		symbol:   market.Symbol(symbol),
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
	// provide the symbol in the context, any producers that care will use it
	ctx = context.WithValue(ctx, market.ContextKeySymbol, p.symbol)
	return NewCsvProducer(res.Body).Produce(ctx, out)
}
