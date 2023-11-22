package csv

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketProducer = (*csvProducer)(nil)

type csvProducer struct {
	data io.ReadCloser
	cols int
}

func NewCsvProducer(data io.ReadCloser) *csvProducer {
	return &csvProducer{
		data: data,
	}
}

// Note: always return an error to stop the producer
func (p *csvProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	reader := bufio.NewReader(p.data)
	header, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	p.cols = len(strings.Split(strings.Trim(header, "\n"), ","))

	for {
		row, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		splts := strings.Split(strings.Trim(row, "\n"), ",")
		if len(splts) != p.cols {
			return io.ErrUnexpectedEOF
		}
		start, err := time.Parse("2006-01-02", splts[0])
		if err != nil {
			return err
		}
		volume, err := strconv.ParseInt(splts[6], 10, 64)
		if err != nil {
			return err
		}

		out <- market.NewMarketLineFromString(start, splts[1], splts[2], splts[3], splts[4], volume)
	}
	return io.ErrClosedPipe
}
