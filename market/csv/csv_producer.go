package csv

import (
	"bufio"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketProducer = (*csvProducer)(nil)

type csvProducer struct {
	fd *os.File
}

func NewCsvProducerFromFile(filePath string) (*csvProducer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &csvProducer{
		fd: file,
	}, nil
}

// Note: always return an error to stop the producer
func (p *csvProducer) Produce(ctx context.Context, out chan market.MarketLine) error {
	reader := bufio.NewReader(p.fd)
	reader.ReadString('\n')
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		splts := strings.Split(strings.Trim(line, "\n"), ",")
		if len(splts) != 7 {
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

func (p *csvProducer) Close() error {
	return p.fd.Close()
}
