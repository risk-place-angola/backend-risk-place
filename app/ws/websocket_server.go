package ws

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var client = make(map[*websocket.Conn]bool)
var lock sync.RWMutex

// WebsocketServer godoc
// @Summary Websocket server
// @Description websocket url ws://host/ws or use authentication ssl wss://host/ws
// @Tags Websocket
// @scheme ws
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /ws [get]
func WebsocketServer(ctx echo.Context) error {

	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return ctx.JSON(400, err)
	}

	defer func(conn *websocket.Conn) {
		lock.Lock()
		delete(client, conn)
		lock.Unlock()
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	lock.Lock()
	client[conn] = true
	lock.Unlock()

	channel := make(chan string)
	go broadcast(channel)

	// remover unreachable code
	// lock.Lock() // lock
	// delete(client, conn)
	// lock.Unlock() // unlock

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return ctx.JSON(400, err)
		}
		channel <- string(msg)
	}
}

func broadcast(c chan string) {
	for {
		msg := <-c
		lock.RLock()
		for conn := range client {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				return
			}
		}
		lock.RUnlock()
	}
}
