package util

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/api/rest/middleware"
	"log"
	"net/http"
	"net/url"
)

type Websocket struct {
	ID                     string
	Conn                   *websocket.Conn
	Send                   chan []byte
	WebsocketClientManager *WebsocketClientManager
}

type WebsocketClientManager struct {
	Clients    map[*Websocket]bool
	Broadcast  chan []byte
	Register   chan *Websocket
	Unregister chan *Websocket
}

func NewWebsocketClientManager() *WebsocketClientManager {
	return &WebsocketClientManager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Websocket),
		Unregister: make(chan *Websocket),
		Clients:    make(map[*Websocket]bool),
	}
}

func (manager *WebsocketClientManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client] = true
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Send)
			}
		case message := <-manager.Broadcast:
			for client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(manager.Clients, client)
				}
			}
		}
	}
}

var uri string

func WebsocketClientDialer(ctx echo.Context) (*websocket.Conn, *echo.HTTPError) {

	uri := Uri(ctx.Request())

	authHeader, err := WebsocketAuthMiddleware(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	conn, _, err := websocket.DefaultDialer.Dial(uri, http.Header{"Authorization": []string{authHeader}})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	log.Println("connected to websocket server")
	return conn, nil
}

func WebsocketAuthMiddleware(ctx echo.Context) (string, error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if ok, err := middleware.IsValidToken(authHeader); !ok || err != nil {
		return "", err
	}
	return authHeader, nil
}

func Uri(r *http.Request) string {
	if r.TLS != nil {
		uri = "wss://" + r.Host + "/ws"
	} else {
		uri = "ws://" + r.Host + "/ws"
	}

	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func (w *Websocket) WebsocketClientWriteMessage(message []byte) {

	defer func(conn *websocket.Conn) {
		w.WebsocketClientManager.Unregister <- w
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
			return
		}
	}(w.Conn)

	err := w.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	w.WebsocketClientManager.Broadcast <- message

}

func (w *Websocket) WebsocketClientReadMessage() ([]byte, error) {

	defer func(conn *websocket.Conn) {
		err := w.Conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
			return
		}
	}(w.Conn)

	for {
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return nil, err
		}
		return message, nil
	}

}

func (w *Websocket) WebsocketServerWriteMessage() {
	defer func() {
		w.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-w.Send:
			if !ok {
				w.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println("WebsocketServerWriteMessage", string(message))
			err := w.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func (w *Websocket) WebsocketServerReadMessage() {
	defer func() {
		w.WebsocketClientManager.Unregister <- w
		w.Conn.Close()
	}()

	for {
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			break
		}
		log.Println("WebsocketServerReadMessage", string(message))
		w.WebsocketClientManager.Broadcast <- message
	}
}