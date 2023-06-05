package ws

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/util"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

	manage := util.NewWebsocketClientManager()
	client := &util.Websocket{
		ID:                     uuid.NewV4().String(),
		Conn:                   conn,
		Send:                   make(chan []byte),
		WebsocketClientManager: manage,
	}

	go manage.Start()

	manage.Register <- client

	go client.WebsocketServerWriteMessage()
	go client.WebsocketServerReadMessage()

	return nil
}
