package place_controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	place_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/place/presenter"
	place_usecase "github.com/risk-place-angola/backend-risk-place/usecase/place"
	uuid "github.com/satori/go.uuid"
)

type PlaceClient struct {
	ID   string
	Conn *websocket.Conn

	Send               chan []byte
	PlaceClientManager *PlaceClientManager
}

type PlaceClientManager struct {
	clients      map[*PlaceClient]bool
	broadcast    chan []byte
	register     chan *PlaceClient
	unregister   chan *PlaceClient
	placeUseCase place_usecase.PlaceUseCase
}

func NewPlaceClientManager(placeUseCase place_usecase.PlaceUseCase) *PlaceClientManager {
	return &PlaceClientManager{
		broadcast:    make(chan []byte),
		register:     make(chan *PlaceClient),
		unregister:   make(chan *PlaceClient),
		clients:      make(map[*PlaceClient]bool),
		placeUseCase: placeUseCase,
	}
}

func (manager *PlaceClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			log.Printf("Client %d connected ID %s ", len(manager.clients), conn.ID)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.Send)
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *PlaceClientManager) Send(message []byte, ignore *PlaceClient) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (manager *PlaceClientManager) PlaceHandler(ctx place_presenter.PlacePresenterCTX) error {
	conn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		return err
	}

	client := &PlaceClient{
		ID:                 uuid.NewV4().String(),
		Conn:               conn,
		Send:               make(chan []byte),
		PlaceClientManager: manager,
	}

	manager.register <- client

	go client.read()
	go client.write()

	return nil
}

func (c *PlaceClient) read() {
	defer func() {
		c.PlaceClientManager.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		c.PlaceClientManager.broadcast <- message
	}
}

func (c *PlaceClient) write() {
	defer func() {
		c.Conn.Close()
	}()

	for range c.Send {
		places, err := c.PlaceClientManager.placeUseCase.FindAllPlace()
		if err != nil {
			if err = c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				log.Println(err)
			}
			return
		}

		if err = c.Conn.WriteJSON(places); err != nil {
			log.Println(err)
		}
	}

	if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		log.Println(err)
	}
}
