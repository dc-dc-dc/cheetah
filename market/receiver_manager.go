package market

import (
	"context"
	"log"
	"sync"
)

type ReceiverManager struct {
	sync.Mutex
	ctx          context.Context
	incomingLine chan MarketLine
	receivers    []MarketReceiver
}

func NewReceiverManager(ctx context.Context) *ReceiverManager {
	rcv := &ReceiverManager{
		ctx:          ctx,
		receivers:    make([]MarketReceiver, 0),
		incomingLine: make(chan MarketLine),
	}
	go rcv.run()
	return rcv
}

func (m *ReceiverManager) run() {
	for {
		select {
		case line := <-m.incomingLine:
			wg := sync.WaitGroup{}

			for index, receiver := range m.receivers {
				wg.Add(1)
				ctx := context.WithValue(m.ctx, "receiver", index)
				var removeSelf func()
				removeSelf = func() {
					m.Lock()
					defer m.Unlock()
					m.receivers = append(m.receivers[:index], m.receivers[index+1:]...)
				}
				go func(ctx context.Context, receiver MarketReceiver, removeSelf func()) {
					if err := receiver.Receive(ctx, line); err != nil {
						log.Printf("receiver err: %v\n", err)
						removeSelf()
					}
					wg.Done()
				}(ctx, receiver, removeSelf)
			}

			wg.Wait()
		}
	}
}

func (m *ReceiverManager) AddReceiver(receiver ...MarketReceiver) {
	m.Lock()
	defer m.Unlock()
	m.receivers = append(m.receivers, receiver...)
}

func (m *ReceiverManager) Receive(ctx context.Context, line MarketLine) error {
	m.incomingLine <- line
	return nil
}
