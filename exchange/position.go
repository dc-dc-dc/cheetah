package exchange

import (
	"time"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

type PositionState int

const (
	PositionStatePending = iota
	PositionStateOpen
	PositionStateClosed
)

type Position struct {
	ID       util.ID
	OpenedAt time.Time
	ClosedAt time.Time

	symbol market.Symbol

	orders []*Order
}

func NewPosition(order *Order) *Position {
	return &Position{
		ID:     util.EnsureID(),
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

func (p Position) ActiveOrders() []*Order {
	orders := make([]*Order, 0)
	for _, o := range p.orders {
		if o.IsOpen() {
			orders = append(orders, o)
		}
	}
	return orders
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

func (p Position) AveragePrice() decimal.Decimal {
	var total decimal.Decimal
	var size int64
	for _, o := range p.orders {
		if o.IsBuy() {
			total = total.Add(o.FilledPrice.Mul(decimal.NewFromInt(o.Filled)))
			size += o.Filled
		} else {
			total = total.Sub(o.FilledPrice.Mul(decimal.NewFromInt(o.Filled)))
			size -= o.Filled
		}
	}
	if size == 0 {
		return decimal.Zero
	}
	return total.Div(decimal.NewFromInt(size))
}

func (p Position) Symbol() market.Symbol {
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
