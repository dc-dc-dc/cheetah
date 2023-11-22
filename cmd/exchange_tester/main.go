package main

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/exchange"
	"github.com/dc-dc-dc/cheetah/exchange/basic"
	"github.com/shopspring/decimal"
)

func main() {
	xgn := basic.NewBasicExchange()
	ctx := context.Background()

	order := exchange.NewMarketBuyOrder("AAPL", decimal.NewFromInt(180), 100)
	if err := xgn.PlaceOrder(ctx, order); err != nil {
		panic(err)
	}

	positions, err := xgn.GetPositions(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range positions {
		fmt.Printf("Position id: %s, size: %d, placed_at: %s, closed_at: %s \n", p.ID.String(), p.Size(), p.OpenedAt, p.ClosedAt)
	}

	if err := xgn.CancelOrder(ctx, order.ID); err != nil {
		panic(err)
	}

	positions, err = xgn.GetPositions(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range positions {
		fmt.Printf("Position id: %s, size: %d, placed_at: %s, closed_at: %s \n", p.ID.String(), p.Size(), p.OpenedAt, p.ClosedAt)
	}

}
