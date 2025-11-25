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
	maxConcurrentBroadcasts             = 100
	broadcastTimeoutSeconds             = 5
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

	broadcastTicker    *time.Ticker
	stopBroadcast      chan bool
	broadcastSemaphore chan struct{}
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
		broadcastSemaphore: make(chan struct{}, maxConcurrentBroadcasts),
	}
}

func (h *Hub) Run() {
	slog.Info("[HUB] WebSocket Hub started successfully")
	go h.startNearbyUsersBroadcast()

	for {
		select {
		case client := <-h.register:
			h.clientsMux.Lock()
			h.clients[client] = true
			slog.Info("[HUB] Client registered", slog.String("user_id", client.UserID), slog.Bool("is_authenticated", client.IsAuthenticated))
			h.clientsMux.Unlock()

		case client := <-h.unregister:
			h.clientsMux.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.closeChannel()
				slog.Info("[HUB] Client unregistered", slog.String("user_id", client.UserID))
			}
			h.clientsMux.Unlock()

		case message := <-h.broadcast:
			h.clientsMux.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					client.closeChannel()
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
	slog.Info("[BROADCAST] Starting nearby users broadcast timer", slog.Int("interval_seconds", nearbyUsersBroadcastIntervalSeconds))
	for {
		select {
		case <-h.broadcastTicker.C:
			h.broadcastNearbyUsersToAll()
		case <-h.stopBroadcast:
			slog.Info("[BROADCAST] Stopping broadcast timer")
			return
		}
	}
}

func (h *Hub) broadcastNearbyUsersToAll() {
	ctx, cancel := context.WithTimeout(context.Background(), broadcastTimeoutSeconds*time.Second)
	defer cancel()

	h.clientsMux.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.clientsMux.RUnlock()

	if len(clients) == 0 {
		return
	}

	slog.Debug("[BROADCAST] Broadcasting nearby users", slog.Int("total_clients", len(clients)))

	for _, client := range clients {
		select {
		case h.broadcastSemaphore <- struct{}{}:
			go func(c *Client) {
				defer func() { <-h.broadcastSemaphore }()
				// Create independent context for each client to avoid cancellation issues
				clientCtx, clientCancel := context.WithTimeout(context.Background(), broadcastTimeoutSeconds*time.Second)
				defer clientCancel()
				h.sendNearbyUsersToClient(clientCtx, c)
			}(client)
		case <-ctx.Done():
			slog.Warn("broadcast timeout reached, skipping remaining clients",
				"remaining", len(clients))
			return
		default:
			slog.Debug("semaphore full, skipping client broadcast",
				"user_id", client.UserID)
		}
	}
}

func (h *Hub) sendNearbyUsersToClient(ctx context.Context, client *Client) {
	if client.lastLat == 0 && client.lastLon == 0 {
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
	}

	const defaultRadius = 5000.0

	users, err := h.nearbyUsersService.GetNearbyUsers(ctx, client.UserID, client.lastLat, client.lastLon, defaultRadius)
	if err != nil {
		slog.Error("failed to get nearby users for broadcast",
			slog.String("user_id", client.UserID),
			slog.Bool("is_authenticated", client.IsAuthenticated),
			slog.Any("error", err))
		return
	}

	slog.Debug("[BROADCAST] Preparing nearby users response",
		slog.String("user_id", client.UserID),
		slog.Int("count", len(users)),
		slog.Float64("client_lat", client.lastLat),
		slog.Float64("client_lon", client.lastLon))

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

		slog.Debug("[BROADCAST] Adding nearby user to response",
			slog.String("anonymous_id", u.AnonymousID),
			slog.Float64("lat", u.Latitude),
			slog.Float64("lon", u.Longitude),
			slog.String("avatar", u.AvatarID),
			slog.String("color", u.Color))
	}

	slog.Debug("[BROADCAST] Sending nearby_users message",
		slog.String("to_user", client.UserID),
		slog.Int("total_users", len(responses)),
		slog.Float64("radius", defaultRadius))

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

		isAnonymous := !c.IsAuthenticated

		slog.Info("updating location via websocket",
			slog.String("user_id", c.UserID),
			slog.Float64("lat", payload.Latitude),
			slog.Float64("lon", payload.Longitude),
			slog.Bool("is_authenticated", c.IsAuthenticated),
			slog.Bool("is_anonymous", isAnonymous))

		// Update Redis location store
		err = h.locationStore.UpdateUserLocation(ctx, c.UserID, payload.Latitude, payload.Longitude)
		if err != nil {
			slog.Error("failed to update location store", slog.Any("error", err))
			c.SendJSON("location_update_failed", map[string]interface{}{
				"status":  "error",
				"message": "Failed to update location in cache",
			})
			return
		}

		// Update PostgreSQL user_locations table
		err = h.nearbyUsersService.UpdateUserLocation(ctx, c.UserID, c.UserID, payload.Latitude, payload.Longitude, payload.Speed, payload.Heading, isAnonymous)
		if err != nil {
			slog.Error("CRITICAL: failed to update user location in nearby service",
				slog.String("user_id", c.UserID),
				slog.Bool("is_anonymous", isAnonymous),
				slog.Any("error", err))
			c.SendJSON("location_update_failed", map[string]interface{}{
				"status":  "error",
				"message": "Failed to update location in database",
				"error":   err.Error(),
			})
			return
		}

		slog.Debug("location updated successfully",
			slog.String("user_id", c.UserID),
			slog.Bool("is_anonymous", isAnonymous))

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
