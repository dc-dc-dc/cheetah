package market_test

import (
	"context"
	"testing"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	ctx := context.Background()
	_, err := market.GetCache(ctx)
	assert.EqualError(t, err, market.ErrNoContextCache.Error())
	_, err = market.GetFromCache[int](ctx, "test")
	assert.EqualError(t, err, market.ErrNoContextCache.Error())

	ctx = market.CreateCache(ctx)
	cache, err := market.GetCache(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	val, err := market.GetFromCache[int](ctx, "test")
	assert.EqualError(t, err, market.ErrNoCacheValue.Error())
	assert.Equal(t, 0, val)

	market.SetCache(ctx, "test", 1)
	_, err = market.GetFromCache[string](ctx, "test")
	assert.EqualError(t, err, market.ErrNoCacheValue.Error())
	val1, err := market.GetFromCache[int](ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, 1, val1)
}
