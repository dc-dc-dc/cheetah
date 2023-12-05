package market_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/stretchr/testify/assert"
)

func TestInterval(t *testing.T) {
	assert.Equal(t, market.Interval1Minute.String(), "1m")
	assert.Equal(t, market.Interval2Minute.String(), "2m")
	assert.Equal(t, market.Interval5Minute.String(), "5m")
	assert.Equal(t, market.Interval15Minute.String(), "15m")
	assert.Equal(t, market.Interval30Minute.String(), "30m")
	assert.Equal(t, market.Interval60Minute.String(), "60m")
	assert.Equal(t, market.Interval90Minute.String(), "90m")
	assert.Equal(t, market.Interval1Hour.String(), "1h")
	assert.Equal(t, market.Interval1Day.String(), "1d")
	assert.Equal(t, market.Interval5Day.String(), "5d")
	assert.Equal(t, market.Interval1Week.String(), "1wk")
	assert.Equal(t, market.Interval1Month.String(), "1mo")
	assert.Equal(t, market.Interval3Month.String(), "3mo")
	assert.Equal(t, market.Interval1Year.String(), "1y")
	assert.Equal(t, market.Interval(100).String(), "")
}
