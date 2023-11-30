package market

import (
	"context"
	"errors"
)

const (
	ContextCache ContextKey = "receiver.cache"
)

var (
	ErrNoContextCache = errors.New("no context cache")
	ErrNoCacheValue   = errors.New("no cache value")
)

type MarketCache map[string]interface{}

func CreateCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextCache, make(MarketCache))
}

func GetCache(ctx context.Context) (MarketCache, error) {
	cache, ok := ctx.Value(ContextCache).(MarketCache)
	if !ok {
		return nil, ErrNoContextCache
	}
	return cache, nil
}

func SetCache(ctx context.Context, key string, value any) {
	cache, err := GetCache(ctx)
	if err != nil {
		return
	}
	cache[key] = value
}

func GetFromCache[T any](ctx context.Context, key string) (T, error) {
	cache, err := GetCache(ctx)
	if err != nil {
		return *new(T), err
	}
	val, ok := cache[key]
	if !ok {
		return *new(T), ErrNoCacheValue
	}
	casted, ok := val.(T)
	if !ok {
		return *new(T), ErrNoCacheValue
	}
	return casted, nil
}
