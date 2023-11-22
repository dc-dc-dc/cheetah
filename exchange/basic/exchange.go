package basic

import (
	"context"
	"time"

	"github.com/dc-dc-dc/cheetah/exchange"
)

type basicExchange struct {
	positions       []*exchange.Position
	activePositions map[string]*exchange.Position
	activeOrders    map[string]*exchange.Order
}

func NewBasicExchange() exchange.Exchange {
	return &basicExchange{
		positions:       make([]*exchange.Position, 0),
		activePositions: make(map[string]*exchange.Position),
		activeOrders:    make(map[string]*exchange.Order),
	}
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
