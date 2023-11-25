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
		indicator.NewMacd(),
		basic.NewBasicReceiver(),
	)
	raw, err := market.SerializeReceivers(receiver)
	if err != nil {
		panic(err)
	}
	// raw := `[{"key":"chained_receiver","raw":"W3sia2V5IjoiaW5kaWNhdG9yLm1pbl9tYXgiLCJyYXciOiJleUozYVc1a2IzY2lPaklzSW0xcGJpSTZkSEoxWlgwPSJ9LHsia2V5IjoiY2hhaW5lZF9yZWNlaXZlciIsInJhdyI6Ilczc2lhMlY1SWpvaWFXNWthV05oZEc5eUxtMXBibDl0WVhnaUxDSnlZWGNpT2lKbGVVb3pZVmMxYTJJelkybFBhazF6U1cweGNHSnBTVFprU0VveFdsZ3dQU0o5WFE9PSJ9LHsia2V5IjoiY2hhaW5lZF9yZWNlaXZlciIsInJhdyI6Ilczc2lhMlY1SWpvaWFXNWthV05oZEc5eUxtMXBibDl0WVhnaUxDSnlZWGNpT2lKbGVVb3pZVmMxYTJJelkybFBhbEZ6U1cweGNHSnBTVFprU0VveFdsZ3dQU0o5WFE9PSJ9XQ=="}]`
	println(string(raw))
	indicator.LoadIndicators()
	res, err := market.DeserializeReceivers([]byte(raw))
	if err != nil {
		panic(err)
	}
	for _, receiver := range res {
		fmt.Printf("%+v\n", receiver)
	}
}
