package exchange

import (
	"time"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

type OrderState int
type OrderType int
type OrderSide int

const (
	OrderSideBuy OrderSide = iota
	OrderSideSell

	OrderStatePending OrderState = iota
	OrderStateOpen
	OrderStateFilled
	OrderStateCanceled

	OrderTypeMarket OrderType = iota
	OrderTypeLimit
	OrderTypeStop
)

type Order struct {
	ID   util.ID
	Side OrderSide
	Type OrderType

	PlacedAt   time.Time
	FilledAt   time.Time
	CanceledAt time.Time

	Price       decimal.Decimal
	FilledPrice decimal.Decimal

	Symbol market.Symbol

	Requested int64
	Filled    int64
}

func NewMarketBuyOrder(symbol market.Symbol, price decimal.Decimal, amount int64) Order {
	return Order{
		ID:        util.EnsureID(),
		Side:      OrderSideBuy,
		Type:      OrderTypeMarket,
		Requested: amount,
		Price:     price,
		Symbol:    symbol,
	}
}

func NewMarketSellOrder(symbol market.Symbol, price decimal.Decimal, amount int64) Order {
	return Order{
		ID:        util.EnsureID(),
		Side:      OrderSideSell,
		Type:      OrderTypeMarket,
		Requested: amount,
		Price:     price,
		Symbol:    symbol,
	}
}

func NewLimitBuyOrder(symbol market.Symbol, price decimal.Decimal, amount int64) Order {
	return Order{
		ID:        util.EnsureID(),
		Side:      OrderSideBuy,
		Type:      OrderTypeLimit,
		Requested: amount,
		Price:     price,
		Symbol:    symbol,
	}
}

func NewLimitSellOrder(symbol market.Symbol, price decimal.Decimal, amount int64) Order {
	return Order{
		ID:        util.EnsureID(),
		Side:      OrderSideSell,
		Type:      OrderTypeLimit,
		Requested: amount,
		Price:     price,
		Symbol:    symbol,
	}
}

func NewStopBuyOrder(symbol market.Symbol, price decimal.Decimal, amount int64) Order {
	return Order{
		ID:        util.EnsureID(),
		Side:      OrderSideBuy,
		Type:      OrderTypeStop,
		Requested: amount,
		Price:     price,
		Symbol:    symbol,
	}
}

func NewStopSellOrder(symbol market.Symbol, price decimal.Decimal, amount int64) Order {
	return Order{
		ID:        util.EnsureID(),
		Side:      OrderSideSell,
		Type:      OrderTypeStop,
		Requested: amount,
		Price:     price,
		Symbol:    symbol,
	}
}

func (o *Order) State() OrderState {
	if !o.CanceledAt.IsZero() {
		return OrderStateCanceled
	}
	if !o.FilledAt.IsZero() {
		return OrderStateFilled
	}
	if !o.PlacedAt.IsZero() {
		return OrderStateOpen
	}
	return OrderStatePending
}

func (o *Order) IsOpen() bool {
	return o.State() == OrderStateOpen
}

func (o *Order) IsMarket() bool {
	return o.Type == OrderTypeMarket
}

func (o *Order) IsLimit() bool {
	return o.Type == OrderTypeLimit
}

func (o *Order) IsBuy() bool {
	return o.Side == OrderSideBuy
}
