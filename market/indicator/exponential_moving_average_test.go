package indicator_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/csv"
	"github.com/dc-dc-dc/cheetah/market/indicator"
	"github.com/shopspring/decimal"
)

func TestExponentialMovingAverage(t *testing.T) {
	file, err := os.Open(testingFileName)
	if err != nil {
		t.Errorf("error opening file: %s  err: %s", testingFileName, err.Error())
		return
	}
	defer file.Close()
	csvReader := csv.NewCsvReader(file)
	header, err := csvReader.Header()
	if err != nil {
		t.Errorf("error reading header: %s", err.Error())
		return
	}
	ema12Index, ok := header["ema12"]
	if !ok {
		t.Errorf("testing data header does not contain ema12")
		return
	}
	ctx := market.CreateCache(context.Background())
	receiver := indicator.NewExponentialMovingAverage(12)
	for {
		line, err := csvReader.NextLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				t.Error(err)
			}
			return
		}
		marketLine, _ := csv.GetMarketLine(header, line)
		if err := receiver.Receive(ctx, *marketLine); err != nil {
			t.Errorf("error receiving market line: %s", err.Error())
			return
		}
		if line[ema12Index] != "" {
			ema, err := market.GetFromCache[decimal.Decimal](ctx, receiver.CacheKey())
			if err != nil {
				t.Errorf("error getting ema from cache: %s", err.Error())
				return
			}
			if !IsInRange(ema, decimal.RequireFromString(line[ema12Index]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("line %d: %s != %s", csvReader.LineNumber(), ema.String(), line[ema12Index])
				return
			}
		}
	}
}
