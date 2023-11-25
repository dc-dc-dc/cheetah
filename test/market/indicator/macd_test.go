package indicator_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/csv"
	"github.com/dc-dc-dc/cheetah/market/indicator"
	"github.com/shopspring/decimal"
)

func TestMacd(t *testing.T) {
	// run the python script and get the results
	// python test/market/indicator/macd_test.py
	fmt.Printf("SKIP_PYTHON = %s\n", os.Getenv("SKIP_PYTHON"))
	symbol := "AAPL"
	timeframe := market.Interval1Day
	startTimeStr := "2019-01-01"
	endTimeStr := "2020-01-01"
	startTime, _ := time.Parse("2006-01-02", startTimeStr)
	endTime, _ := time.Parse("2006-01-02", endTimeStr)
	if os.Getenv("SKIP_PYTHON") == "true" {
		fmt.Printf("running python\n")
		cmd := exec.Command("./macd_test.py")
		cmd.Env = append(cmd.Env, fmt.Sprintf("SYMBOL=%s", symbol))
		cmd.Env = append(cmd.Env, fmt.Sprintf("TIMEFRAME=%s", timeframe))
		cmd.Env = append(cmd.Env, fmt.Sprintf("START_TIME=%s", startTimeStr))
		cmd.Env = append(cmd.Env, fmt.Sprintf("END_TIME=%s", endTimeStr))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Error(err)
		}
	}
	producer := csv.NewYFinanceProducer(symbol, market.Interval1Day, startTime, endTime)
	out := make(chan market.MarketLine)
	ctx := context.Background()

	go func(ctx context.Context, out chan market.MarketLine) {
		for {
			if err := producer.Produce(ctx, out); err != nil {
				if !errors.Is(err, io.ErrClosedPipe) {
					fmt.Printf("[producer] err: %v\n", err)
				}
				return
			}
		}
	}(ctx, out)
	fd, err := os.Open("macd_test_py.csv")
	if err != nil {
		t.Error(err)
	}
	defer fd.Close()
	csvReader := csv.NewCsvReader(fd)
	header, err := csvReader.Header()
	if err != nil {
		t.Error(err)
	}

	var getStats = func() (decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
		line, err := csvReader.NextLine()
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, err
		}
		var ema12, ema26, macd decimal.Decimal
		if ema12, err = decimal.NewFromString(line[header["ema12"]]); err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, err
		}
		if ema26, err = decimal.NewFromString(line[header["ema26"]]); err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, err
		}
		if macd, err = decimal.NewFromString(line[header["macd"]]); err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, err
		}
		return ema12, ema26, macd, nil
	}

	receiver := market.NewChainedReceiver(
		indicator.NewMacd(),
		market.NewFunctionalReceiver(func(ctx context.Context, line market.MarketLine) error {
			ema12, ema26, macd, err := getStats()
			if err != nil {
				return err
			}
			cmpEma12, err := market.GetFromCache[decimal.Decimal](ctx, indicator.ExponentialMovingAverageCacheKey(12))
			if err != nil {
				return err
			}
			cmpEma26, err := market.GetFromCache[decimal.Decimal](ctx, indicator.ExponentialMovingAverageCacheKey(26))
			if err != nil {
				return err
			}

			cmpMacd, err := market.GetFromCache[decimal.Decimal](ctx, indicator.MacdCacheKey())
			if err != nil {
				return err
			}
			fmt.Printf("ema12: %s, cmp: %s\n", ema12.String(), cmpEma12.String())
			if false {
				fmt.Printf("ema26: %s, cmp: %s\n", ema26.String(), cmpEma26.String())
				fmt.Printf("macd: %s, cmp: %s\n", macd.String(), cmpMacd.String())
			}
			return nil
		}),
	)

	for {
		select {
		case <-ctx.Done():
			return
		case line := <-out:
			if market.IsMarketNull(line) {
				return
			}
			if err := receiver.Receive(ctx, line); err != nil {
				t.Error(err)
			}
		}
	}
}
