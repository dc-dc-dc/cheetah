package market

import "context"

type MarketReceiver interface {
	Receive(context.Context, MarketLine) error
}

type CachableReceiver interface {
	CacheKey() string
	MarketReceiver
}

type FunctionalReceiver struct {
	call func(context.Context, MarketLine) error
}

func NewFunctionalReceiver(call func(context.Context, MarketLine) error) MarketReceiver {
	return &FunctionalReceiver{call: call}
}

func (r *FunctionalReceiver) Receive(ctx context.Context, line MarketLine) error {
	return r.call(ctx, line)
}
