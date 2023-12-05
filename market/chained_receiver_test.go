package market_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/basic"
	"github.com/dc-dc-dc/cheetah/market/indicator"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/stretchr/testify/assert"
)

func TestChainedReceiver(t *testing.T) {
	cr := market.NewChainedReceiver()
	assert.Equal(t, cr.PrefixKey(), "chained_receiver")
	assert.Empty(t, cr.Receivers())
	assert.Equal(t, cr.String(), "ChainedReceiver{receivers=[]}")

	gen, ok := market.GetSerializableReceiverGenerator(cr.PrefixKey())
	assert.True(t, ok)
	assert.NotNil(t, gen)
	crFromGen := gen()
	assert.NotNil(t, crFromGen)
	assert.IsType(t, cr, crFromGen)

	raw, err := json.Marshal(cr)
	assert.NoError(t, err)
	assert.Equal(t, "[]", string(raw))
	err = json.Unmarshal([]byte("{\"window\": \"testing\"}"), crFromGen)
	assert.Error(t, err)
	err = json.Unmarshal(raw, crFromGen)
	assert.NoError(t, err)
}

func TestChainedReceiverReceive(t *testing.T) {
	cr := market.NewChainedReceiver(basic.NewBasicReceiver())
	err := cr.Receive(context.Background(), market.MarketLine{})
	assert.NoError(t, err)
	cr = market.NewChainedReceiver(basic.NewErrorReceiver(0))
	err = cr.Receive(context.Background(), market.MarketLine{})
	assert.EqualError(t, err, "receiver error")
}

func TestDedupReceivers(t *testing.T) {

	cr := market.DedupReceivers([]market.MarketReceiver{indicator.NewMinIndicator(1), indicator.NewMinIndicator(1), basic.NewBasicReceiver()}, util.NewSet[string]())
	assert.Len(t, cr, 2)

	cr = market.DedupReceivers([]market.MarketReceiver{indicator.NewMinIndicator(1), market.NewChainedReceiver(indicator.NewMinIndicator(1))}, util.NewSet[string]())
	assert.Len(t, cr, 1)

	cr = market.DedupReceivers([]market.MarketReceiver{indicator.NewMinIndicator(1), market.NewChainedReceiver(indicator.NewMinIndicator(2))}, util.NewSet[string]())
	assert.Len(t, cr, 2)
}
