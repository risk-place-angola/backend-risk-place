package handler

import (
	"log/slog"
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
)

type RiskHandler struct {
	app *application.Application
}

func NewRiskHandler(app *application.Application) *RiskHandler {
	return &RiskHandler{
		app: app,
	}
}

// ListRiskTypes godoc.
// @Summary List all risk types.
// @Description Retrieve all available risk types with their default radius.
// @Tags risks
// @Accept json
// @Produce json
// @Success 200 {object} dto.RiskTypesListResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /risks/types [get].
func (h *RiskHandler) ListRiskTypes(w http.ResponseWriter, r *http.Request) {
	riskTypes, err := h.app.RiskUseCase.ListRiskTypes(r.Context())
	if err != nil {
		slog.Error("failed to list risk types", "error", err)
		util.Error(w, "failed to retrieve risk types", http.StatusInternalServerError)
		return
	}

	util.Response(w, riskTypes, http.StatusOK)
}

// ListRiskTopics godoc.
// @Summary List risk topics.
// @Description Retrieve risk topics, optionally filtered by risk_type_id query parameter.
// @Tags risks
// @Accept json
// @Produce json
// @Param risk_type_id query string false "Filter by risk type ID (UUID)"
// @Success 200 {object} dto.RiskTopicsListResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /risks/topics [get].
func (h *RiskHandler) ListRiskTopics(w http.ResponseWriter, r *http.Request) {
	riskTypeID := r.URL.Query().Get("risk_type_id")

	var riskTypeIDPtr *string
	if riskTypeID != "" {
		riskTypeIDPtr = &riskTypeID
	}

	riskTopics, err := h.app.RiskUseCase.ListRiskTopics(r.Context(), riskTypeIDPtr)
	if err != nil {
		slog.Error("failed to list risk topics", "error", err)
		util.Error(w, "failed to retrieve risk topics", http.StatusInternalServerError)
		return
	}

	util.Response(w, riskTopics, http.StatusOK)
}
