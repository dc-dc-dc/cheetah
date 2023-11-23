package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dc-dc-dc/cheetah/exchange"
	"github.com/dc-dc-dc/cheetah/exchange/basic"
	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/csv"
	"github.com/shopspring/decimal"
)

func main() {
	xgn := basic.NewBasicExchange()
	ctx := context.Background()

	order := exchange.NewMarketBuyOrder("AAPL", decimal.NewFromInt(160), 100)
	if err := xgn.PlaceOrder(ctx, order); err != nil {
		panic(err)
	}

	positions, err := xgn.GetPositions(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range positions {
		printPosition(p)
	}

	producer := csv.NewYFinanceProducer("AAPL", market.Interval1Day, time.Now().Add(-time.Hour*24*365), time.Now())
	out := make(chan market.MarketLine)
	go func(out chan market.MarketLine) {
		if err := producer.Produce(ctx, out); err != nil {
			close(out)
		}
	}(out)
	for line := range out {
		ctx = context.WithValue(ctx, market.ContextKeySymbol, market.Symbol("AAPL"))
		if err := xgn.Receive(ctx, line); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}

	positions, err = xgn.GetPositions(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range positions {
		printPosition(p)
	}
}

func printPosition(p exchange.Position) {

	fmt.Printf("Position id: %s, size: %d, placed_at: %s", p.ID.String(), p.Size(), p.OpenedAt)
	if p.State() == exchange.PositionStateClosed {
		fmt.Printf(", closed_at: %s", p.ClosedAt)
	}
	if p.State() == exchange.PositionStateOpen {
		fmt.Printf(", average_price: %s", p.AveragePrice())
	}

	fmt.Print("\n")
}
