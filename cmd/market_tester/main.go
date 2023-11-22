package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/basic"
	"github.com/dc-dc-dc/cheetah/market/csv"
)

func fakeProducer() market.MarketProducer {
	lines := []market.MarketLine{
		market.EnsureMarketLineFromString(time.Now().Add(1*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.EnsureMarketLineFromString(time.Now().Add(2*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.EnsureMarketLineFromString(time.Now().Add(3*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.EnsureMarketLineFromString(time.Now().Add(4*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.EnsureMarketLineFromString(time.Now().Add(5*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.EnsureMarketLineFromString(time.Now().Add(6*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.EnsureMarketLineFromString(time.Now().Add(7*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
	}
	return basic.NewBasicProducer(lines, 1)
}

func main() {
	fmt.Printf("Cheetah a market tool potentially.\n")
	out := make(chan market.MarketLine)
	ctx := context.Background()

	// producer := fakeProducer()
	// file, err := os.Open("data/apple.csv")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }
	// defer file.Close()
	// producer := csv.NewCSVProducer(file)
	producer := csv.NewYFinanceProducer("AAPL", csv.YFinanceInterval1Day, time.Now().Add(-1*365*24*time.Hour), time.Now())
	rcvMgr := market.NewReceiverManager(ctx)
	rcvMgr.AddReceiver(basic.NewBasicReceiver(), basic.NewCountReceiver())

	// For testing purposes...
	// Create a producer
	//
	go func(ctx context.Context, out chan market.MarketLine) {
		for {
			if err := producer.Produce(ctx, out); err != nil {
				fmt.Printf("producer err: %v\n", err)
				return
			}
		}
	}(ctx, out)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case line := <-out:
				if err := rcvMgr.Receive(ctx, line); err != nil {
					fmt.Printf("listener err: %v\n", err)
					return
				}
			}
		}
	}()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	<-sig
}
