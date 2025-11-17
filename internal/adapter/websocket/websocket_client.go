package websocket

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID          string // Can be user_id (UUID) or device_id
	IsAuthenticated bool   // true if JWT authenticated, false if anonymous
	Conn            *websocket.Conn
	Send            chan []byte
	Hub             *Hub
}

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.Hub.unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		c.Hub.handleIncomingMessage(ctx, c, message)
	}
}

func (c *Client) WritePump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			slog.Error("error closing connection", slog.Any("error", err))
		}
	}(c.Conn)
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			slog.Error("error writing message", slog.Any("error", err))
			return
		}
	}

	// Channel was closed, send close message
	err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		slog.Error("error writing close message", slog.Any("error", err))
	}
}

func (c *Client) SendJSON(event string, data interface{}) {
	msg, err := json.Marshal(Message{Event: event, Data: data})
	if err != nil {
		slog.Error("error marshaling message", slog.Any("error", err))
		return
	}
	c.Send <- msg
}
