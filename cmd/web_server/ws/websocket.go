package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/gorilla/websocket"
)

type MessageWrapper struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Handler func(ctx context.Context, payload interface{}, id util.ID) error

type WebSocketManager struct {
	upgrader websocket.Upgrader
	clients  sync.Map
	handlers map[string]Handler
	ctx      context.Context
	cancel   func()
}

func NewWebSocketManager(ctx context.Context) *WebSocketManager {
	ctx, cancel := context.WithCancel(ctx)
	return &WebSocketManager{
		ctx:    ctx,
		cancel: cancel,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		},
		clients:  sync.Map{},
		handlers: make(map[string]Handler),
	}
}

func (ws *WebSocketManager) Close() {
	ws.cancel()
	ws.clients.Range(func(key any, val any) bool {
		fmt.Printf("closing client id=%s\n", key)
		val.(*WebSocketClient).conn.Close()
		return true
	})
}

func (ws *WebSocketManager) HttpHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[error] err: %s\n", err)
	}
	wsClient := NewWebSocketClient(conn)
	ws.clients.Store(wsClient.ID(), wsClient)
	go ws.listenClient(wsClient)
}

func (ws *WebSocketManager) RegisterHandler(key string, handler Handler) {
	// TODO: Multiple handlers for the same thing ????
	ws.handlers[key] = handler
}

func (ws *WebSocketManager) SendMessage(id util.ID, msg MessageWrapper) error {
	client, ok := ws.clients.Load(id)
	if !ok {
		return fmt.Errorf("no client with id=%s", id)
	}
	return client.(*WebSocketClient).SendMessage(msg)
}

func (ws *WebSocketManager) listenClient(client *WebSocketClient) {
	ctx, cancel := context.WithCancel(ws.ctx)
	defer cancel()
	for {
		_, data, err := client.conn.ReadMessage()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				fmt.Printf("[web-socket] got error client=%s err=%s\n", client.ID(), err)
			}
			ws.clients.Delete(client.ID())
			return
		}
		var msg MessageWrapper
		if err := json.Unmarshal(data, &msg); err != nil {
			fmt.Printf("[web-socket] got error client=%s err=%s\n", client.ID(), err)
			continue
		}
		handler, ok := ws.handlers[msg.Type]
		if !ok {
			fmt.Printf("[web-socket] error unknown message type=%s\n", msg.Type)
			continue
		}
		go func(ctx context.Context, msg MessageWrapper, id util.ID) {
			if err := handler(ctx, msg.Payload, id); err != nil {
				fmt.Printf("[web-socket] got error running handler for type=%s err=%s\n", msg.Type, err)
			}
		}(ctx, msg, client.ID())
	}
}
