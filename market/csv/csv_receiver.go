package csv

import (
	"context"
	"io"
	"strconv"

	"github.com/dc-dc-dc/cheetah/market"
)

func NewCsvReceiver(dst io.Writer) market.MarketReceiver {
	writer := NewCsvWriter(&dst)
	return market.NewFunctionalReceiver(func(ctx context.Context, line market.MarketLine) error {
		elements := make([]string, 0, 5)
		elements = append(elements, line.Start.Format("2006-01-02 15:04:05"))
		elements = append(elements, line.Open.String())
		elements = append(elements, line.High.String())
		elements = append(elements, line.Low.String())
		elements = append(elements, line.Close.String())
		elements = append(elements, strconv.FormatInt(line.Volume, 10))
		return writer.Write(elements)
	})
}
