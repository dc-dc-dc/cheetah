package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/basic"
)

func main() {
	fmt.Printf("Cheetah a market tool potentially.\n")
	lines := []market.MarketLine{
		market.NewMarketLineFromString(time.Now().Add(1*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.NewMarketLineFromString(time.Now().Add(2*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.NewMarketLineFromString(time.Now().Add(3*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.NewMarketLineFromString(time.Now().Add(4*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.NewMarketLineFromString(time.Now().Add(5*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.NewMarketLineFromString(time.Now().Add(6*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
		market.NewMarketLineFromString(time.Now().Add(7*time.Minute), "8.0", "10.0", "7.0", "7.5", 1),
	}
	basicProducer := basic.NewBasicProducer(lines, 1)

	out := make(chan market.MarketLine)
	ctx := context.Background()

	rcvMgr := market.NewReceiverManager(ctx)
	rcvMgr.AddReceiver(basic.NewBasicReceiver())

	// For testing purposes...
	// Create a producer
	//
	// Have receivers(s)
	//  - Bot (making trades)
	go func(ctx context.Context, out chan market.MarketLine) {
		for {
			if err := basicProducer.Produce(ctx, out); err != nil {
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
