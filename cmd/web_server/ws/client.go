package ws

import (
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	id   util.ID
	conn *websocket.Conn
}

func NewWebSocketClient(conn *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{
		id:   util.EnsureID(),
		conn: conn,
	}
}

func (client *WebSocketClient) SendMessage(msg MessageWrapper) error {
	return client.conn.WriteJSON(msg)
}

func (client *WebSocketClient) ID() util.ID {
	return client.id
}
