package websocket

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
)

type WSHandler struct {
	Hub            *Hub
	AuthMiddleware middleware.AuthMiddleware
}

func NewWSHandler(hub *Hub, authMiddleware middleware.AuthMiddleware) *WSHandler {
	return &WSHandler{
		Hub:            hub,
		AuthMiddleware: authMiddleware,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket godoc
// @Summary Handle WebSocket connections for alerts
// @Description Upgrade HTTP connection to WebSocket for real-time alerts
// @Tags websocket
// @Security BearerAuth
// @Success 101 {string} string "Switching Protocols"
// @Failure 401 {object} util.ErrorResponse
// @Router /ws/alerts [get]
func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userIDStr, err := h.AuthMiddleware.ValidateJWTFromRequest(r)
	if err != nil {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to WebSocket", slog.Any("error", err))
		util.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	client := &Client{
		UserID: userIDStr,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.Hub,
	}

	slog.Info("websocket client connected", slog.String("user_id", userIDStr))

	h.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
