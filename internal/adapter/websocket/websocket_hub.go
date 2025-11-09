package websocket

import (
	"context"
	"encoding/json"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"log"
	"log/slog"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	locationStore port.LocationStore
	geoService    port.GeolocationService
}

func NewHub(locationStore port.LocationStore, geoService port.GeolocationService) *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan []byte),
		locationStore: locationStore,
		geoService:    geoService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) handleIncomingMessage(c *Client, raw []byte) {
	var msg Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Println("invalid message:", err)
		return
	}

	switch msg.Event {
	case "update_location":
		var payload UpdateLocationPayload
		b, _ := json.Marshal(msg.Data)
		_ = json.Unmarshal(b, &payload)

		err := h.locationStore.UpdateUserLocation(context.Background(), c.UserID, payload.Latitude, payload.Longitude)
		if err != nil {
			log.Printf("failed to update location: %v", err)
			return
		}

		c.SendJSON("location_updated", map[string]interface{}{
			"status": "ok",
		})
	}
}

func (h *Hub) BroadcastAlert(ctx context.Context, alertID string, message string, lat, lon, radius float64) {
	userIDs, err := h.locationStore.FindUsersInRadius(ctx, lat, lon, radius)
	if err != nil {
		slog.Error("error finding nearby users", "error", err)
		return
	}

	slog.Info("broadcasting alert", "alert_id", alertID, "user_count", len(userIDs))

	for client := range h.clients {
		for _, id := range userIDs {
			if client.UserID == id {
				client.SendJSON("new_alert", AlertNotification{
					AlertID:   alertID,
					Message:   message,
					Latitude:  lat,
					Longitude: lon,
					Radius:    radius,
				})
			}
		}
	}
}

func (h *Hub) BroadcastReport(ctx context.Context, reportID, message string, lat, lon, radius float64) {
	userIDs, _ := h.locationStore.FindUsersInRadius(ctx, lat, lon, radius)
	for client := range h.clients {
		for _, id := range userIDs {
			if client.UserID == id {
				client.SendJSON("report_created", ReportNotification{
					ReportID:  reportID,
					Message:   message,
					Latitude:  lat,
					Longitude: lon,
				})
			}
		}
	}
}

func (h *Hub) NotifyUser(userID string, event string, data interface{}) {
	for client := range h.clients {
		if client.UserID == userID {
			client.SendJSON(event, data)
		}
	}
}

func (h *Hub) Clients() map[*Client]bool {
	return h.clients
}
