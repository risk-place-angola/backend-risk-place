package risk_controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	risk_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/risk/presenter"
	risk_usecase "github.com/risk-place-angola/backend-risk-place/usecase/risk"
	uuid "github.com/satori/go.uuid"
)

type RiskClient struct {
	ID   string
	Conn *websocket.Conn

	Send              chan []byte
	RiskClientManager *RiskClientManager
}

type RiskClientManager struct {
	clients     map[*RiskClient]bool
	broadcast   chan []byte
	register    chan *RiskClient
	unregister  chan *RiskClient
	riskUseCase risk_usecase.RiskUseCase
}

func NewRiskClientManager(riskUseCase risk_usecase.RiskUseCase) *RiskClientManager {
	return &RiskClientManager{
		broadcast:   make(chan []byte),
		register:    make(chan *RiskClient),
		unregister:  make(chan *RiskClient),
		clients:     make(map[*RiskClient]bool),
		riskUseCase: riskUseCase,
	}
}

func (manager *RiskClientManager) Start() {
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

func (manager *RiskClientManager) Send(message []byte, ignore *RiskClient) {
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

func (manager *RiskClientManager) RiskHandler(ctx risk_presenter.RiskPresenterCTX) error {
	conn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		return err
	}

	client := &RiskClient{
		ID:                uuid.NewV4().String(),
		Conn:              conn,
		Send:              make(chan []byte),
		RiskClientManager: manager,
	}

	manager.register <- client

	go client.read()
	go client.write()

	return nil
}

func (c *RiskClient) read() {
	defer func() {
		c.RiskClientManager.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		c.RiskClientManager.broadcast <- message
	}
}

func (c *RiskClient) write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case _, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			risks, err := c.RiskClientManager.riskUseCase.FindAllRisk()
			if err != nil {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteJSON(risks)

		}
	}
}
