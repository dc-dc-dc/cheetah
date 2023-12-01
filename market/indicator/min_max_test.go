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

func TestMinMax(t *testing.T) {
	maxIndicator := indicator.NewMaxIndicator(20)
	minIndicator := indicator.NewMinIndicator(20)
	assert.Equal(t, "indicator.min_max", maxIndicator.PrefixKey())
	assert.Equal(t, "indicator.min_max.20.max", maxIndicator.CacheKey())
	assert.Equal(t, "indicator.min_max.20.min", minIndicator.CacheKey())
	assert.Equal(t, "MinMaxIndicator{window=20, min=false}", maxIndicator.String())
	assert.Equal(t, "MinMaxIndicator{window=20, min=true}", minIndicator.String())

	raw, err := json.Marshal(maxIndicator)
	assert.NoError(t, err)
	assert.Equal(t, "{\"window\":20,\"min\":false}", string(raw))

	gen, ok := market.GetSerializableReceiverGenerator(maxIndicator.PrefixKey())
	assert.True(t, ok)
	assert.NotNil(t, gen)
	receiver := gen()
	assert.NotNil(t, receiver)
	assert.IsType(t, maxIndicator, receiver)
	err = json.Unmarshal([]byte("{\"window\": \"testing\"}"), receiver)
	assert.Error(t, err)
	err = json.Unmarshal(raw, receiver)
	assert.NoError(t, err)

	assert.Equal(t, maxIndicator, receiver)
}

func TestMinMaxMovingAverageCalc(t *testing.T) {
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
	maxIndicator20Index, ok := header["max20"]
	if !ok {
		t.Errorf("testing data header does not contain maxIndicator20")
		return
	}
	minIndicator20Index, ok := header["min20"]
	if !ok {
		t.Errorf("testing data header does not contain minIndicator20")
		return
	}
	ctx := market.CreateCache(context.Background())
	maxReceiver := indicator.NewMaxIndicator(20)
	minReceiver := indicator.NewMinIndicator(20)
	for {
		line, err := csvReader.NextLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				t.Error(err)
			}
			return
		}
		marketLine, _ := csv.GetMarketLine(header, line)
		if err := maxReceiver.Receive(ctx, *marketLine); err != nil {
			t.Errorf("error receiving market line: %s", err.Error())
			return
		}
		if err := minReceiver.Receive(ctx, *marketLine); err != nil {
			t.Errorf("error receiving market line: %s", err.Error())
			return
		}
		if line[maxIndicator20Index] != "" {
			maxIndicator, err := indicator.GetMaxFromCache(ctx, 20)
			if err != nil {
				t.Errorf("error getting maxIndicator from cache: %s", err.Error())
				return
			}
			if !IsInRange(maxIndicator, decimal.RequireFromString(line[maxIndicator20Index]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("line %d: %s != %s", csvReader.LineNumber(), maxIndicator.String(), line[maxIndicator20Index])
				return
			}
			minIndicator, err := indicator.GetMinFromCache(ctx, 20)
			if err != nil {
				t.Errorf("error getting minIndicator from cache: %s", err.Error())
				return
			}
			if !IsInRange(minIndicator, decimal.RequireFromString(line[minIndicator20Index]), decimal.NewFromFloat32(0.0001)) {
				t.Errorf("line %d: %s != %s", csvReader.LineNumber(), minIndicator.String(), line[minIndicator20Index])
				return
			}

		}
	}
}
