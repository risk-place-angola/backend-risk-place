package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

type DangerZoneHandler struct {
	app *application.Application
}

func NewDangerZoneHandler(app *application.Application) *DangerZoneHandler {
	return &DangerZoneHandler{app: app}
}

// GetDangerZonesNearby godoc.
// @Summary Get nearby danger zones.
// @Description Retrieves danger zones near a specific location based on incident density and risk score.
// @Tags danger-zones
// @Accept json
// @Produce json
// @Security OptionalAuth
// @Param X-Device-Id header string false "Device ID for anonymous users"
// @Param request body dto.GetDangerZonesRequest true "Location and radius"
// @Success 200 {object} dto.GetDangerZonesResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /danger-zones/nearby [post].
func (h *DangerZoneHandler) GetDangerZonesNearby(w http.ResponseWriter, r *http.Request) {
	var req dto.GetDangerZonesRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode danger zones request", "error", err)
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.app.DangerZoneUseCase.GetDangerZonesNearby(ctx, &req)
	if err != nil {
		slog.Error("failed to get danger zones", "error", err)
		util.Error(w, "failed to retrieve danger zones", http.StatusInternalServerError)
		return
	}

	util.Response(w, response, http.StatusOK)
}
