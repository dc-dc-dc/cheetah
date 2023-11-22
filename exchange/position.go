package exchange

import (
	"time"
)

type PositionState int

const (
	PositionStatePending = iota
	PositionStateOpen
	PositionStateClosed
)

type Position struct {
	ID       ID
	OpenedAt time.Time
	ClosedAt time.Time

	symbol string

	orders []*Order
}

func NewPosition(order *Order) *Position {
	return &Position{
		ID:     EnsureID(),
		symbol: order.Symbol,
		orders: []*Order{order},
	}
}

func (p Position) HasOpen() bool {
	for _, o := range p.orders {
		if o.IsOpen() {
			return true
		}
	}
	return false
}

func (p Position) Size() int64 {
	var size int64
	for _, o := range p.orders {
		if o.IsBuy() {
			size += o.Filled
		} else {
			size -= o.Filled
		}
	}
	return size
}

func (p Position) Symbol() string {
	return p.symbol
}

func (p Position) AddOrder(o *Order) {
	if !p.ClosedAt.IsZero() {
		return
	}
	p.orders = append(p.orders, o)
}

func (p Position) State() PositionState {
	if !p.ClosedAt.IsZero() {
		return PositionStateClosed
	}
	if !p.OpenedAt.IsZero() {
		return PositionStateOpen
	}
	return PositionStatePending
}
