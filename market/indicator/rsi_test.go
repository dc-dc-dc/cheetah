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

func TestRsi(t *testing.T) {
	rsigen, ok := market.GetSerializableReceiverGenerator(indicator.RsiCacheKey())
	assert.True(t, ok)
	assert.NotNil(t, rsigen)
	rsiFromGen := rsigen()
	assert.NotNil(t, rsiFromGen)
	rsi := indicator.NewRsi()
	assert.NotNil(t, rsi)

	// when nothing is in cache this should not return an error
	ctx := market.CreateCache(context.Background())
	err := rsi.Receive(ctx, market.MarketLine{})
	assert.NoError(t, err)
}

func TestRsiCalc(t *testing.T) {
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
	rsiIndex, ok := header["rsi14"]
	if !ok {
		t.Errorf("testing data header does not contain macd")
		return
	}
	ctx := market.CreateCache(context.Background())
	receiver := indicator.NewRsi()
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
		if line[rsiIndex] != "" {
			rsi, err := indicator.GetRsiFromCache(ctx)
			if err != nil {
				t.Errorf("error getting ema from cache: %s", err.Error())
				return
			}
			if !IsInRange(rsi, decimal.RequireFromString(line[rsiIndex]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("rsi line %d: got: %s expected: %s", csvReader.LineNumber(), rsi.String(), line[rsiIndex])
			}

		}
	}
}
