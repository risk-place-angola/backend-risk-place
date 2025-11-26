package websocket

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
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
	settingsChecker    service.SettingsChecker

	broadcastTicker    *time.Ticker
	stopBroadcast      chan bool
	broadcastSemaphore chan struct{}
}

func NewHub(locationStore port.LocationStore, geoService port.GeolocationService, nearbyUsersService port.NearbyUsersService, settingsChecker service.SettingsChecker) *Hub {
	return &Hub{
		clients:            make(map[*Client]bool),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		broadcast:          make(chan []byte),
		locationStore:      locationStore,
		geoService:         geoService,
		nearbyUsersService: nearbyUsersService,
		settingsChecker:    settingsChecker,
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

func (h *Hub) BroadcastAlert(ctx context.Context, alertID string, message string, lat, lon, radius float64, severity string) {
	userIDs, err := h.locationStore.FindUsersInRadius(ctx, lat, lon, radius)
	if err != nil {
		slog.Error("error finding nearby users", "error", err)
		return
	}

	slog.Debug("found users in radius for alert broadcast", "alert_id", alertID, "potential_users", len(userIDs))

	notifiedCount := 0
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	for client := range h.clients {
		for _, userIDStr := range userIDs {
			if client.UserID != userIDStr {
				continue
			}

			userUUID, err := uuid.Parse(userIDStr)
			if err != nil {
				slog.Debug("invalid user UUID format, skipping", "user_id", userIDStr)
				continue
			}

			deviceID := ""
			if !client.IsAuthenticated {
				deviceID = userIDStr
				userUUID = uuid.Nil
			}

			if !h.settingsChecker.CanReceiveNotifications(ctx, userUUID, deviceID) {
				continue
			}

			distanceMeters := h.calculateDistance(lat, lon, client.lastLat, client.lastLon)

			isHighRiskTime := h.settingsChecker.IsInHighRiskTime(ctx, userUUID, deviceID)
			if !isHighRiskTime {
				if !h.settingsChecker.CanReceiveAlerts(ctx, userUUID, deviceID, severity, int(distanceMeters)) {
					continue
				}
			} else {
				if distanceMeters > maxSearchRadiusMeters {
					continue
				}
				slog.Debug("time-based boost applied", "user_id", userIDStr, "severity", severity, "is_high_risk_time", true)
			}

			client.SendJSON("new_alert", AlertNotification{
				AlertID:   alertID,
				Message:   message,
				Latitude:  lat,
				Longitude: lon,
				Radius:    radius,
			})
			notifiedCount++
			break
		}
	}

	slog.Info("alert broadcast completed", "alert_id", alertID, "notified_users", notifiedCount, "potential_users", len(userIDs))
}

func (h *Hub) BroadcastReport(ctx context.Context, reportID, message string, lat, lon, radius float64, isVerified bool) {
	userIDs, err := h.locationStore.FindUsersInRadius(ctx, lat, lon, radius)
	if err != nil {
		slog.Error("error finding nearby users for report", "error", err)
		return
	}

	slog.Debug("found users in radius for report broadcast", "report_id", reportID, "potential_users", len(userIDs))

	notifiedCount := 0
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	for client := range h.clients {
		for _, userIDStr := range userIDs {
			if client.UserID != userIDStr {
				continue
			}

			userUUID, err := uuid.Parse(userIDStr)
			if err != nil {
				slog.Debug("invalid user UUID format for report, skipping", "user_id", userIDStr)
				continue
			}

			deviceID := ""
			if !client.IsAuthenticated {
				deviceID = userIDStr
				userUUID = uuid.Nil
			}

			if !h.settingsChecker.CanReceiveNotifications(ctx, userUUID, deviceID) {
				continue
			}

			distanceMeters := h.calculateDistance(lat, lon, client.lastLat, client.lastLon)

			if !h.settingsChecker.CanReceiveReports(ctx, userUUID, deviceID, isVerified, int(distanceMeters)) {
				continue
			}

			client.SendJSON("report_created", ReportNotification{
				ReportID:  reportID,
				Message:   message,
				Latitude:  lat,
				Longitude: lon,
			})
			notifiedCount++
			break
		}
	}

	slog.Info("report broadcast completed", "report_id", reportID, "notified_users", notifiedCount, "potential_users", len(userIDs))
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

func (h *Hub) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	if lat2 == 0 && lon2 == 0 {
		return 0
	}

	const (
		earthRadiusMeters = 6371000.0
		degreesToRadians  = math.Pi / 180.0
		halfDivisor       = 2.0
	)

	dLat := (lat2 - lat1) * degreesToRadians
	dLon := (lon2 - lon1) * degreesToRadians

	halfDLat := dLat / halfDivisor
	halfDLon := dLon / halfDivisor

	a := math.Sin(halfDLat)*math.Sin(halfDLat) +
		math.Cos(lat1*degreesToRadians)*math.Cos(lat2*degreesToRadians)*
			math.Sin(halfDLon)*math.Sin(halfDLon)

	c := halfDivisor * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}
