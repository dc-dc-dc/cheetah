package basic

import (
	"context"
	"io"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketProducer = (*basicProducer)(nil)

type basicProducer struct {
	sleepTimeSeconds time.Duration
	line             []market.MarketLine
}

func NewBasicProducer(lines []market.MarketLine, sleepTimeSeconds int64) *basicProducer {
	return &basicProducer{
		line:             lines,
		sleepTimeSeconds: time.Duration(sleepTimeSeconds),
	}
}

func (p *basicProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	for _, line := range p.line {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- line:
			time.Sleep(p.sleepTimeSeconds * time.Second)
		}
	}
	return io.ErrClosedPipe
}
