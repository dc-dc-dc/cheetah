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
	"github.com/stretchr/testify/assert"
)

func TestMacd(t *testing.T) {
	macdgen, ok := market.GetSerializableReceiverGenerator(indicator.MacdCacheKey())
	assert.True(t, ok)
	assert.NotNil(t, macdgen)
	macdFromGen := macdgen()
	assert.NotNil(t, macdFromGen)
	macd := indicator.NewMacd()
	assert.NotNil(t, macd)

	// when nothing is in cache this should not return an error
	ctx := market.CreateCache(context.Background())
	err := macd.Receive(ctx, market.MarketLine{})
	assert.NoError(t, err)
}

func TestMacdCalc(t *testing.T) {
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
	macdIndex, ok := header["macd"]
	if !ok {
		t.Errorf("testing data header does not contain macd")
		return
	}
	signalIndex, ok := header["macd.signal"]
	if !ok {
		t.Errorf("testing data header does not contain macd.signal")
		return
	}
	histogramIndex, ok := header["macd.histogram"]
	if !ok {
		t.Errorf("testing data header does not contain macd.histogram")
		return
	}
	ctx := market.CreateCache(context.Background())
	receiver := indicator.NewMacd()
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
		if line[macdIndex] != "" {
			macd, err := indicator.GetMacdFromCache(ctx)
			if err != nil {
				t.Errorf("error getting ema from cache: %s", err.Error())
				return
			}
			if !IsInRange(macd.Macd, decimal.RequireFromString(line[macdIndex]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("macd line %d: got: %s expected: %s", csvReader.LineNumber(), macd.Macd.String(), line[macdIndex])
			}
			if !IsInRange(macd.Signal, decimal.RequireFromString(line[signalIndex]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("signal line %d: got: %s expected: %s", csvReader.LineNumber(), macd.Signal.String(), line[signalIndex])
			}
			if !IsInRange(macd.Histogram, decimal.RequireFromString(line[histogramIndex]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("histogram line %d: got: %s != expected: %s", csvReader.LineNumber(), macd.Histogram.String(), line[histogramIndex])
			}

		}
	}
}
