package main

import (
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/basic"
	"github.com/dc-dc-dc/cheetah/market/indicator"
)

func main() {
	receiver := market.NewChainedReceiver(
		indicator.NewMinIndicator(2),
		market.NewChainedReceiver(indicator.NewMinIndicator(3)),
		market.NewChainedReceiver(indicator.NewMinIndicator(4)),
		indicator.NewExponentialMovingAverage(5),
		basic.NewBasicReceiver(),
	)
	raw, err := market.SerializeReceivers(receiver)
	if err != nil {
		panic(err)
	}
	// res := indicator.MinMaxIndicator{}
	res, err := market.DeserializeReceivers(raw)
	if err != nil {
		panic(err)
	}
	for _, receiver := range res {
		fmt.Printf("%+v\n", receiver)
	}
}
