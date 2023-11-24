package basic

import (
	"context"
	"fmt"
	"time"

	"github.com/dc-dc-dc/cheetah/exchange"
	"github.com/dc-dc-dc/cheetah/market"
)

var _ exchange.Exchange = (*basicExchange)(nil)
var _ market.MarketReceiver = (*basicExchange)(nil)

type basicExchange struct {
	positions       []*exchange.Position
	activePositions map[market.Symbol]*exchange.Position
	activeOrders    map[string]*exchange.Order
}

func NewBasicExchange() *basicExchange {
	return &basicExchange{
		positions:       make([]*exchange.Position, 0),
		activePositions: make(map[market.Symbol]*exchange.Position),
		activeOrders:    make(map[string]*exchange.Order),
	}
}

// for emulation purposes, this will mark active orders as complete if the fill price is less than the Close price
func (b *basicExchange) Receive(ctx context.Context, line market.MarketLine) error {
	symbol, ok := ctx.Value(market.ContextKeySymbol).(market.Symbol)
	if !ok {
		return market.ErrSymbolNotFound
	}
	// get the active positions
	activePosition, ok := b.activePositions[symbol]
	if !ok {
		// if there is no active position, then there is nothing to do
		return nil
	}
	// get the active orders
	activeOrders := activePosition.ActiveOrders()
	for _, order := range activeOrders {
		if order.IsBuy() {
			switch order.Type {
			case exchange.OrderTypeMarket:
				{
					order.Filled = order.Requested
					order.FilledPrice = line.Close
					order.FilledAt = time.Now()
					break
				}
			case exchange.OrderTypeLimit:
				{
					if order.Price.GreaterThanOrEqual(line.Close) {
						order.Filled = order.Requested
						order.FilledPrice = order.Price
						order.FilledAt = time.Now()
						break
					}
				}
			default:
				{
					fmt.Printf("unknown order type: %v\n", order.Type)
				}
			}

		} else {
			switch order.Type {
			case exchange.OrderTypeMarket:
				{
					order.Filled = order.Requested
					order.FilledPrice = line.Close
					order.FilledAt = time.Now()
					break
				}
			case exchange.OrderTypeLimit:
				{
					if order.Price.LessThanOrEqual(line.Close) {
						order.Filled = order.Requested
						order.FilledPrice = order.Price
						order.FilledAt = time.Now()
						break
					}
				}
			default:
				{
					fmt.Printf("unknown order type: %v\n", order.Type)
				}
			}
		}
	}
	return nil
}

func (b *basicExchange) GetPositions(ctx context.Context) ([]exchange.Position, error) {
	res := make([]exchange.Position, len(b.positions))
	for i, p := range b.positions {
		res[i] = *p
	}
	return res, nil
}

func (b *basicExchange) PlaceOrder(ctx context.Context, order exchange.Order) error {
	b.activeOrders[order.ID.String()] = &order
	activePosition, ok := b.activePositions[order.Symbol]
	if !ok {
		activePosition = exchange.NewPosition(&order)
		b.activePositions[order.Symbol] = activePosition
		activePosition.OpenedAt = time.Now()
		b.positions = append(b.positions, activePosition)
	}
	activePosition.AddOrder(&order)
	order.PlacedAt = time.Now()
	return nil
}

func (b *basicExchange) CancelOrder(ctx context.Context, orderId exchange.ID) error {
	order, ok := b.activeOrders[orderId.String()]
	if !ok {
		return exchange.ErrOrderNotFound
	}
	order.CanceledAt = time.Now()
	position := b.activePositions[order.Symbol]
	if position.Size() == 0 {
		position.ClosedAt = time.Now()
		delete(b.activePositions, order.Symbol)
	}
	return nil
}
