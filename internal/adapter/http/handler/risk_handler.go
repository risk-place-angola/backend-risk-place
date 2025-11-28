package handler

import (
	"encoding/json"
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

// GetRiskType godoc
// @Summary Get a risk type by ID
// @Description Retrieve a specific risk type by its UUID
// @Tags risks
// @Accept json
// @Produce json
// @Param id path string true "Risk Type ID (UUID)"
// @Success 200 {object} dto.RiskTypeResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /risks/types/{id} [get]
func (h *RiskHandler) GetRiskType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		util.Error(w, "risk type ID is required", http.StatusBadRequest)
		return
	}

	riskType, err := h.app.RiskUseCase.GetRiskType(r.Context(), id)
	if err != nil {
		slog.Error("failed to get risk type", "id", id, "error", err)
		util.Error(w, "risk type not found", http.StatusNotFound)
		return
	}

	util.Response(w, riskType, http.StatusOK)
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

// GetRiskTopic godoc
// @Summary Get a risk topic by ID
// @Description Retrieve a specific risk topic by its UUID
// @Tags risks
// @Accept json
// @Produce json
// @Param id path string true "Risk Topic ID (UUID)"
// @Success 200 {object} dto.RiskTopicResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /risks/topics/{id} [get]
func (h *RiskHandler) GetRiskTopic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		util.Error(w, "risk topic ID is required", http.StatusBadRequest)
		return
	}

	riskTopic, err := h.app.RiskUseCase.GetRiskTopic(r.Context(), id)
	if err != nil {
		slog.Error("failed to get risk topic", "id", id, "error", err)
		util.Error(w, "risk topic not found", http.StatusNotFound)
		return
	}

	util.Response(w, riskTopic, http.StatusOK)
}

func (h *RiskHandler) UpdateRiskTypeIsEnabled(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		util.Error(w, "risk type ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		IsEnabled bool `json:"is_enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.app.RiskUseCase.UpdateRiskTypeIsEnabled(r.Context(), id, req.IsEnabled); err != nil {
		slog.Error("failed to update risk type enabled status", "risk_type_id", id, "error", err)
		util.Error(w, "failed to update risk type", http.StatusInternalServerError)
		return
	}

	slog.Info("risk type status updated", "risk_type_id", id, "is_enabled", req.IsEnabled)
	util.Response(w, map[string]interface{}{
		"message":    "risk type status updated successfully",
		"is_enabled": req.IsEnabled,
	}, http.StatusOK)
}
