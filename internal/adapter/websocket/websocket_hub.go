package websocket

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
)

const (
	nearbyUsersBroadcastIntervalSeconds = 3
	defaultSearchRadiusMeters           = 5000.0
	maxSearchRadiusMeters               = 10000.0
)

type Hub struct {
	clients    map[*Client]bool
	clientsMux sync.RWMutex
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	locationStore      port.LocationStore
	geoService         port.GeolocationService
	nearbyUsersService port.NearbyUsersService

	broadcastTicker *time.Ticker
	stopBroadcast   chan bool
}

func NewHub(locationStore port.LocationStore, geoService port.GeolocationService, nearbyUsersService port.NearbyUsersService) *Hub {
	return &Hub{
		clients:            make(map[*Client]bool),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		broadcast:          make(chan []byte),
		locationStore:      locationStore,
		geoService:         geoService,
		nearbyUsersService: nearbyUsersService,
		broadcastTicker:    time.NewTicker(nearbyUsersBroadcastIntervalSeconds * time.Second),
		stopBroadcast:      make(chan bool),
	}
}

func (h *Hub) Run() {
	go h.startNearbyUsersBroadcast()

	for {
		select {
		case client := <-h.register:
			h.clientsMux.Lock()
			h.clients[client] = true
			h.clientsMux.Unlock()

		case client := <-h.unregister:
			h.clientsMux.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.clientsMux.Unlock()

		case message := <-h.broadcast:
			h.clientsMux.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					h.clientsMux.RUnlock()
					h.clientsMux.Lock()
					delete(h.clients, client)
					h.clientsMux.Unlock()
					h.clientsMux.RLock()
				}
			}
			h.clientsMux.RUnlock()

		case <-h.stopBroadcast:
			h.broadcastTicker.Stop()
			return
		}
	}
}

func (h *Hub) startNearbyUsersBroadcast() {
	for {
		select {
		case <-h.broadcastTicker.C:
			h.broadcastNearbyUsersToAll()
		case <-h.stopBroadcast:
			return
		}
	}
}

func (h *Hub) broadcastNearbyUsersToAll() {
	ctx := context.Background()
	h.clientsMux.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.clientsMux.RUnlock()

	for _, client := range clients {
		go h.sendNearbyUsersToClient(ctx, client)
	}
}

func (h *Hub) sendNearbyUsersToClient(ctx context.Context, client *Client) {
	if client.lastLat == 0 && client.lastLon == 0 {
		return
	}

	const defaultRadius = 5000.0
	users, err := h.nearbyUsersService.GetNearbyUsers(ctx, client.UserID, client.lastLat, client.lastLon, defaultRadius)
	if err != nil {
		slog.Error("failed to get nearby users for broadcast", "user_id", client.UserID, "error", err)
		return
	}

	responses := make([]NearbyUserResponse, len(users))
	for i, u := range users {
		responses[i] = NearbyUserResponse{
			UserID:    u.AnonymousID,
			Latitude:  u.Latitude,
			Longitude: u.Longitude,
			AvatarID:  u.AvatarID,
			Color:     u.Color,
			Speed:     u.Speed,
			Heading:   u.Heading,
		}
	}

	client.SendJSON("nearby_users", NearbyUsersData{
		Users:      responses,
		Radius:     defaultRadius,
		TotalCount: len(responses),
	})
}

func (h *Hub) handleIncomingMessage(ctx context.Context, c *Client, raw []byte) {
	var msg Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Println("invalid message:", err)
		return
	}

	switch msg.Event {
	case "update_location":
		var payload UpdateLocationPayload
		b, err := json.Marshal(msg.Data)
		if err != nil {
			log.Printf("failed to marshal payload: %v", err)
			return
		}
		err = json.Unmarshal(b, &payload)
		if err != nil {
			log.Printf("failed to unmarshal payload: %v", err)
			return
		}

		c.lastLat = payload.Latitude
		c.lastLon = payload.Longitude

		err = h.locationStore.UpdateUserLocation(ctx, c.UserID, payload.Latitude, payload.Longitude)
		if err != nil {
			log.Printf("failed to update location: %v", err)
			return
		}

		isAnonymous := !c.IsAuthenticated
		err = h.nearbyUsersService.UpdateUserLocation(ctx, c.UserID, c.UserID, payload.Latitude, payload.Longitude, payload.Speed, payload.Heading, isAnonymous)
		if err != nil {
			log.Printf("failed to update user location: %v", err)
		}

		c.SendJSON("location_updated", map[string]interface{}{"status": "ok"})

	case "get_nearby_users":
		var payload UpdateLocationPayload
		b, err := json.Marshal(msg.Data)
		if err != nil {
			log.Printf("failed to marshal payload: %v", err)
			return
		}
		err = json.Unmarshal(b, &payload)
		if err != nil {
			log.Printf("failed to unmarshal payload: %v", err)
			return
		}

		radius := payload.Radius
		if radius == 0 {
			radius = defaultSearchRadiusMeters
		}
		if radius > maxSearchRadiusMeters {
			radius = maxSearchRadiusMeters
		}

		users, err := h.nearbyUsersService.GetNearbyUsers(ctx, c.UserID, payload.Latitude, payload.Longitude, radius)
		if err != nil {
			log.Printf("failed to get nearby users: %v", err)
			return
		}

		responses := make([]NearbyUserResponse, len(users))
		for i, u := range users {
			responses[i] = NearbyUserResponse{
				UserID:    u.AnonymousID,
				Latitude:  u.Latitude,
				Longitude: u.Longitude,
				AvatarID:  u.AvatarID,
				Color:     u.Color,
				Speed:     u.Speed,
				Heading:   u.Heading,
			}
		}

		c.SendJSON("nearby_users", NearbyUsersData{
			Users:      responses,
			Radius:     radius,
			TotalCount: len(responses),
		})

	default:
		log.Printf("unknown event type: %s", msg.Event)
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
