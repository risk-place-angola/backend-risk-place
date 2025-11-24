package websocket

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

const (
	clientSendBufferSize = 256
)

type Client struct {
	UserID          string
	IsAuthenticated bool
	Conn            *websocket.Conn
	Send            chan []byte
	Hub             *Hub
	lastLat         float64
	lastLon         float64
	closed          atomic.Bool
	closeMux        sync.Once
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
	defer func() {
		c.closeMux.Do(func() {
			c.closed.Store(true)
			close(c.Send)
		})
		if err := c.Conn.Close(); err != nil {
			slog.Error("error closing connection", slog.Any("error", err))
		}
	}()

	for message := range c.Send {
		if c.closed.Load() {
			return
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			slog.Error("error writing message", slog.Any("error", err))
			return
		}
	}

	if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		slog.Error("error writing close message", slog.Any("error", err))
	}
}

func (c *Client) SendJSON(event string, data interface{}) {
	if c.closed.Load() {
		return
	}

	msg, err := json.Marshal(Message{Event: event, Data: data})
	if err != nil {
		slog.Error("error marshaling message", slog.Any("error", err))
		return
	}

	defer func() {
		if r := recover(); r != nil {
			slog.Warn("recovered from panic when sending message", "panic", r, "user_id", c.UserID)
		}
	}()

	select {
	case c.Send <- msg:
	default:
		slog.Warn("client send channel is full or closed, dropping message",
			"user_id", c.UserID,
			"event", event)
	}
}
