package csv

import (
	"context"
	"fmt"
	"io"
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

func (p *yFinanceProducer) Fetch(ctx context.Context) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=%s&events=history", p.symbol, p.start.Unix(), p.end.Unix(), p.interval.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return res.Body, nil
}

func (p *yFinanceProducer) String() string {
	return fmt.Sprintf("yFinanceProducer{symbol: %s, interval: %s, start: %s, end: %s}", p.symbol, p.interval, p.start.Format("2006-01-02 15:04:05"), p.end.Format("2006-01-02 15:04:05"))
}

func (p *yFinanceProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	res, err := p.Fetch(ctx)
	if err != nil {
		return err
	}
	defer res.Close()
	ctx = context.WithValue(ctx, market.ContextKeySymbol, p.symbol)
	return NewCsvProducer(res).Produce(ctx, out)
}
