package websocket

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
)

type WSHandler struct {
	Hub                *Hub
	AuthMiddleware     middleware.AuthMiddleware
	OptionalMiddleware *middleware.OptionalAuthMiddleware
	upgrader           websocket.Upgrader
}

func NewWSHandler(hub *Hub, authMiddleware middleware.AuthMiddleware, optionalMiddleware *middleware.OptionalAuthMiddleware) *WSHandler {
	return &WSHandler{
		Hub:                hub,
		AuthMiddleware:     authMiddleware,
		OptionalMiddleware: optionalMiddleware,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// HandleWebSocket godoc
// @Summary Handle WebSocket connections for alerts
// @Description Upgrade HTTP connection to WebSocket for real-time alerts. Supports both authenticated (JWT) and anonymous (device_id) connections
// @Tags websocket
// @Security BearerAuth
// @Param X-Device-ID header string false "Device ID for anonymous users"
// @Success 101 {string} string "Switching Protocols"
// @Failure 401 {object} util.ErrorResponse
// @Router /ws/alerts [get]
func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Try to extract either JWT or device_id
	identifier, isAuthenticated, err := h.OptionalMiddleware.ExtractIdentifier(r)
	if err != nil {
		slog.Error("failed to extract identifier", slog.Any("error", err))
		util.Error(w, "unauthorized: JWT or X-Device-ID header required", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to WebSocket", slog.Any("error", err))
		util.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	client := &Client{
		UserID:          identifier,
		IsAuthenticated: isAuthenticated,
		Conn:            conn,
		Send:            make(chan []byte, clientSendBufferSize),
		Hub:             h.Hub,
	}
	client.closed.Store(false)

	clientType := "anonymous"
	if isAuthenticated {
		clientType = "authenticated"
	}
	slog.Info("websocket client connected",
		slog.String("identifier", identifier),
		slog.String("type", clientType))

	h.Hub.register <- client

	go client.WritePump()
	// WebSocket connections need a long-lived context independent of the HTTP request
	// r.Context() is canceled after the upgrade, so we use context.Background()
	//nolint:contextcheck // WebSocket requires independent context after HTTP upgrade
	go client.ReadPump(context.Background())
}
