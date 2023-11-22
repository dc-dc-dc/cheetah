package market

import (
	"context"
)

type MarketProducer interface {
	Produce(context.Context, chan MarketLine) error
}
