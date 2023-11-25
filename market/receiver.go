package market

import (
	"context"
	"encoding/json"
)

type MarketReceiveFunc func(context.Context, MarketLine) error

type MarketReceiver interface {
	Receive(context.Context, MarketLine) error
}

type SerializableReceiver interface {
	PrefixKey() string
	MarketReceiver
}

type CachableReceiver interface {
	CacheKey() string
	SerializableReceiver
}

type SerializableDataReceiver interface {
	json.Marshaler
	json.Unmarshaler
	SerializableReceiver
}

type FunctionalReceiver struct {
	call MarketReceiveFunc
}

func NewFunctionalReceiver(call MarketReceiveFunc) MarketReceiver {
	return &FunctionalReceiver{call: call}
}

func (r *FunctionalReceiver) Receive(ctx context.Context, line MarketLine) error {
	return r.call(ctx, line)
}

type CachableFunctionalReceiver struct {
	key  string
	call MarketReceiveFunc
}

func NewCachableFunctionalReceiver(key string, call MarketReceiveFunc) CachableReceiver {
	return &CachableFunctionalReceiver{key: key, call: call}
}

func (r *CachableFunctionalReceiver) PrefixKey() string {
	return r.key

}

func (r *CachableFunctionalReceiver) CacheKey() string {
	return r.key
}

func (r *CachableFunctionalReceiver) Receive(ctx context.Context, line MarketLine) error {
	return r.call(ctx, line)
}

func (r *CachableFunctionalReceiver) String() string {
	return r.key
}
