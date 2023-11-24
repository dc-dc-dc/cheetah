package market

import "context"

type MarketReceiver interface {
	Receive(context.Context, MarketLine) error
}

type CachableReceiver interface {
	CacheKey() string
	MarketReceiver
}
