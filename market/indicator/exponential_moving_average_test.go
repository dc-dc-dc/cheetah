package indicator_test

import (
	"context"
	"encoding/json"
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

func TestExponentialMovingAverage(t *testing.T) {
	sma := indicator.NewExponentialMovingAverage(20)
	assert.Equal(t, "indicator.moving_average_exponential", sma.PrefixKey())
	assert.Equal(t, "indicator.moving_average_exponential.20", sma.CacheKey())
	assert.Equal(t, "ExponentialMovingAverage{window=20}", sma.String())
	raw, err := json.Marshal(sma)
	assert.NoError(t, err)
	assert.Equal(t, "{\"window\":20}", string(raw))
	gen, ok := market.GetSerializableReceiverGenerator(sma.PrefixKey())
	assert.True(t, ok)
	assert.NotNil(t, gen)
	receiver := gen()
	assert.NotNil(t, receiver)
	assert.IsType(t, sma, receiver)
	err = json.Unmarshal([]byte("{\"window\": \"testing\"}"), receiver)
	assert.Error(t, err)
	err = json.Unmarshal(raw, receiver)
	assert.NoError(t, err)

	assert.Equal(t, sma, receiver)
}

func TestExponentialMovingAverageCalc(t *testing.T) {
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
	window := 12
	receiver := indicator.NewExponentialMovingAverage(window)
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
			ema, err := indicator.GetExponentialMovingAverageFromCache(ctx, window)
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
