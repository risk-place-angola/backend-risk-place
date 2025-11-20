package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	httputil "github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
)

const (
	defaultRadiusMeters    = 5000.0
	maxRadiusMeters        = 10000.0
	cleanupIntervalSeconds = 30
	rateLimitSeconds       = 3
)

type NearbyUsersHandler struct {
	nearbyUsersService port.NearbyUsersService
	lastRequests       map[string]time.Time
}

func NewNearbyUsersHandler(nearbyUsersService port.NearbyUsersService) *NearbyUsersHandler {
	return &NearbyUsersHandler{
		nearbyUsersService: nearbyUsersService,
		lastRequests:       make(map[string]time.Time),
	}
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
}

type GetNearbyUsersRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
	Radius    float64 `json:"radius"`
}

type NearbyUsersResponse struct {
	Users      []UserLocation `json:"users"`
	Radius     float64        `json:"radius"`
	TotalCount int            `json:"total_count"`
}

type UserLocation struct {
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	AvatarID  string  `json:"avatar_id"`
	Color     string  `json:"color"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
}

func (h *NearbyUsersHandler) checkRateLimit(userID string) bool {
	now := time.Now()
	if last, ok := h.lastRequests[userID]; ok {
		if now.Sub(last) < rateLimitSeconds*time.Second {
			return false
		}
	}
	h.lastRequests[userID] = now
	return true
}

// UpdateLocation godoc
// @Summary Update user location for nearby users feature
// @Description Update current user location to be visible on the map for nearby users
// @Tags nearby-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Device-ID header string false "Device ID for anonymous users"
// @Param request body UpdateLocationRequest true "Location data with optional speed and heading"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /api/v1/users/location [post]
func (h *NearbyUsersHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		httputil.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = userID
	}

	var req UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", slog.Any("error", err))
		httputil.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	isAnonymous := r.Context().Value(httputil.IsAuthenticatedKey)
	anonymous := false
	if isAnonymous != nil {
		isAuth, ok := isAnonymous.(bool)
		if !ok {
			slog.Error("failed to convert isAuthenticated to bool")
			httputil.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		anonymous = !isAuth
	}

	if err := h.nearbyUsersService.UpdateUserLocation(r.Context(), userID, deviceID, req.Latitude, req.Longitude, req.Speed, req.Heading, anonymous); err != nil {
		slog.Error("failed to update user location", slog.Any("error", err))
		httputil.Error(w, "failed to update location", http.StatusInternalServerError)
		return
	}

	httputil.Response(w, map[string]string{"status": "ok"}, http.StatusOK)
}

// GetNearbyUsers godoc
// @Summary Get nearby users on the map
// @Description Retrieve list of nearby users with their anonymous avatars within specified radius
// @Tags nearby-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Device-ID header string false "Device ID for anonymous users"
// @Param request body GetNearbyUsersRequest true "Location and radius for searching nearby users"
// @Success 200 {object} NearbyUsersResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 429 {object} util.ErrorResponse "Rate limit exceeded (max 1 request per 3 seconds)"
// @Failure 500 {object} util.ErrorResponse
// @Router /api/v1/users/nearby [post]
func (h *NearbyUsersHandler) GetNearbyUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		httputil.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !h.checkRateLimit(userID) {
		httputil.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var req GetNearbyUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", slog.Any("error", err))
		httputil.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Radius == 0 {
		req.Radius = defaultRadiusMeters
	}

	if req.Radius > maxRadiusMeters {
		req.Radius = maxRadiusMeters
	}

	users, err := h.nearbyUsersService.GetNearbyUsers(r.Context(), userID, req.Latitude, req.Longitude, req.Radius)
	if err != nil {
		slog.Error("failed to get nearby users", slog.Any("error", err))
		httputil.Error(w, "failed to get nearby users", http.StatusInternalServerError)
		return
	}

	response := NearbyUsersResponse{
		Users:      make([]UserLocation, len(users)),
		Radius:     req.Radius,
		TotalCount: len(users),
	}

	for i, u := range users {
		response.Users[i] = UserLocation{
			UserID:    u.AnonymousID,
			Latitude:  u.Latitude,
			Longitude: u.Longitude,
			AvatarID:  u.AvatarID,
			Color:     u.Color,
			Speed:     u.Speed,
			Heading:   u.Heading,
		}
	}

	httputil.Response(w, response, http.StatusOK)
}

func StartCleanupJob(ctx context.Context, service port.NearbyUsersService) {
	ticker := time.NewTicker(cleanupIntervalSeconds * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := service.CleanupStaleLocations(ctx); err != nil {
					slog.Error("failed to cleanup stale locations", slog.Any("error", err))
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
