package basic

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market"
)

var _ market.MarketReceiver = (*basicReceiver)(nil)
var _ market.MarketReceiver = (*errorReceiver)(nil)
var _ market.MarketReceiver = (*countReceiver)(nil)

type countReceiver struct {
	count int64
}

func NewCountReceiver() *countReceiver {
	return &countReceiver{
		count: 0,
	}
}

func (r *countReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	r.count++
	if r.count%10 == 0 {
		fmt.Printf("handled %d lines\n", r.count)
	}
	return nil
}

type errorReceiver struct {
	current   int
	explodeOn int
}

func NewErrorReceiver(explodesOn int) *errorReceiver {
	return &errorReceiver{
		explodeOn: explodesOn,
	}
}

func (r *errorReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	r.current++
	logLine(ctx, line)
	if r.current == r.explodeOn {
		return fmt.Errorf("receiver error")
	}
	return nil
}

type basicReceiver struct {
	keys []string
}

func NewBasicReceiver(keys ...string) *basicReceiver {
	return &basicReceiver{
		keys: keys,
	}
}

func (r *basicReceiver) Receive(ctx context.Context, line market.MarketLine) error {
	logLine(ctx, line, r.keys...)
	return nil
}

func logLine(ctx context.Context, line market.MarketLine, keys ...string) {
	index, ok := ctx.Value(market.ContextReceiverIndex).(int)
	var str string
	if ok {
		str += fmt.Sprintf("[receiver-%d]", index)
	}
	str += fmt.Sprintf(" line: %v", line)
	if len(keys) > 0 {
		cache, ok := ctx.Value(market.ContextCache).(map[string]interface{})
		if ok {
			for _, key := range keys {
				value, ok := cache[key]
				if ok {
					str += fmt.Sprintf(" %s: %v", key, value)
				}
			}
		}
	}
	fmt.Print(str + "\n")
}
