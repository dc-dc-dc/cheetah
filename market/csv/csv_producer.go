package csv

import (
	"context"
	"fmt"
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

func GetMarketLine(header CsvHeader, splts []string) (*market.MarketLine, error) {
	if len(splts) == 0 {
		return nil, fmt.Errorf("empty line")
	}
	var timeParse string = "2006-01-02 15:04:05-07:00"
	if !strings.Contains(splts[header["date"]], " ") {
		timeParse = "2006-01-02"
	}
	start, err := time.Parse(timeParse, splts[header["date"]])
	if err != nil {
		return nil, err
	}
	volume, err := strconv.ParseInt(splts[header["volume"]], 10, 64)
	if err != nil {
		return nil, err
	}

	return market.NewMarketLineFromString(
		start,
		splts[header["open"]],
		splts[header["high"]],
		splts[header["low"]],
		splts[header["close"]],
		volume,
	)
}

// Note: always return an error to stop the producer
func (p *csvProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	header, err := p.reader.Header()
	if err != nil {
		return err
	}
	var line *market.MarketLine
	for {
		splts, err := p.reader.NextLine()
		if err != nil {
			break
		}
		line, err = GetMarketLine(header, splts)
		if err != nil {
			break
		}

		if err != nil {
			break
		}
		out <- *line
	}
	close(out)
	return err
}
