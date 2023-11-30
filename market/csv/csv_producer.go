package csv

import (
	"context"
	"io"
	"strconv"
	"strings"
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
	for {
		splts, err := p.reader.NextLine()
		if err != nil {
			break
		}
		var timeParse string = "2006-01-02 15:04:05-07:00"
		if !strings.Contains(splts[header["date"]], " ") {
			timeParse = "2006-01-02"
		}
		start, err := time.Parse(timeParse, splts[header["date"]])
		if err != nil {
			break
		}
		volume, err := strconv.ParseInt(splts[header["volume"]], 10, 64)
		if err != nil {
			break
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
			break
		}
		out <- *line
	}
	close(out)
	return err
}
