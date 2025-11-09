package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

type AlertHandler struct {
	alertUseCase *application.Application
}

func NewAlertHandler(alertUseCase *application.Application) *AlertHandler {
	return &AlertHandler{
		alertUseCase: alertUseCase,
	}
}

// CreateAlert godoc
// @Summary Create a new alert
// @Description Create a new alert
// @Tags alerts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param alert body dto.Alert true "Alert"
// @Success 201 {object} dto.Alert
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /alerts [post]
func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.Alert
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	req.UserID = userIDStr

	err := h.alertUseCase.AlertUseCase.TriggerAlert(r.Context(), req)
	if err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{"status": "alert triggered"}, http.StatusCreated)
}
