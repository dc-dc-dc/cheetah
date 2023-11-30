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
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

var (
	testingFileName = util.GetEnv("TESTING_DATA_FILE", "./testing_data.csv")
)

func IsInRange(test, cmp, delta decimal.Decimal) bool {
	return test.GreaterThanOrEqual(cmp.Sub(delta)) && test.LessThanOrEqual(cmp.Add(delta))
}

func TestSimpleMovingAverage(t *testing.T) {
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
	sma20Index, ok := header["sma20"]
	if !ok {
		t.Errorf("testing data header does not contain sma20")
		return
	}
	ctx := market.CreateCache(context.Background())
	receiver := indicator.NewSimpleMovingAverage(20)
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
		if line[sma20Index] != "" {
			sma, err := market.GetFromCache[decimal.Decimal](ctx, receiver.CacheKey())
			if err != nil {
				t.Errorf("error getting sma from cache: %s", err.Error())
				return
			}
			if !IsInRange(sma, decimal.RequireFromString(line[sma20Index]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("line %d: %s != %s", csvReader.LineNumber(), sma.String(), line[sma20Index])
				return
			}
		}
	}
}
