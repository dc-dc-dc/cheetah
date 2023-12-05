package market

import (
	"context"
	"errors"
	"sync"
)

const (
	ContextCache ContextKey = "receiver.cache"
)

var (
	ErrNoContextCache = errors.New("no context cache")
	ErrNoCacheValue   = errors.New("no cache value")
)

func CreateCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextCache, &sync.Map{})
}

func GetCache(ctx context.Context) (*sync.Map, error) {
	cache, ok := ctx.Value(ContextCache).(*sync.Map)
	if !ok {
		return nil, ErrNoContextCache
	}
	return cache, nil
}

func SetCache(ctx context.Context, key string, value any) {
	if cache, err := GetCache(ctx); err == nil {
		cache.Store(key, value)

	}
}

func GetFromCache[T any](ctx context.Context, key string) (T, error) {
	cache, err := GetCache(ctx)
	if err != nil {
		return *new(T), err
	}
	val, ok := cache.Load(key)
	if !ok {
		return *new(T), ErrNoCacheValue
	}
	casted, ok := val.(T)
	if !ok {
		return *new(T), ErrNoCacheValue
	}
	return casted, nil
}
