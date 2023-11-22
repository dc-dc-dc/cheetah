package market

import "context"

type MarketReceiver interface {
	Receive(context.Context, MarketLine) error
}
