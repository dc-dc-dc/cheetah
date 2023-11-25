package csv

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketProducer = (*csvProducer)(nil)

type csvProducer struct {
	reader *CsvReader
}

func NewCsvProducer(data io.ReadCloser) *csvProducer {
	return &csvProducer{
		reader: NewCsvReader(data),
	}
}

// Note: always return an error to stop the producer
func (p *csvProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	header, err := p.reader.Header()
	if err != nil {
		return err
	}
	fmt.Printf("header: %+v\n", header)
	for {
		splts, err := p.reader.NextLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		start, err := time.Parse("2006-01-02", splts[header["date"]])
		if err != nil {
			return err
		}
		volume, err := strconv.ParseInt(splts[header["volume"]], 10, 64)
		if err != nil {
			return err
		}

		line, err := market.NewMarketLineFromString(
			start,
			splts[header["open"]],
			splts[header["high"]],
			splts[header["low"]],
			splts[header["close"]],
			volume,
		)
		if err != nil {
			return err
		}
		out <- *line
	}
	close(out)
	return io.ErrClosedPipe
}
