package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

type ReportHandler struct {
	reportUseCase *application.Application
}

func NewReportHandler(reportUseCase *application.Application) *ReportHandler {
	return &ReportHandler{
		reportUseCase: reportUseCase,
	}
}

// Create godoc
// @Summary Create a new report
// @Description Create a new report
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param report body dto.ReportCreate true "Report to create"
// @Success 201 {object} dto.ReportDTO
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /reports [post]
func (h *ReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.ReportCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.UserID = userIDStr

	report, err := h.reportUseCase.ReportUseCase.Create(r.Context(), req)
	if err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dto.ReportToDTO(report)
	util.Response(w, res, http.StatusCreated)
}

// ListNearby godoc
// @Summary List nearby reports
// @Description List reports near the specified location
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lat query string true "Latitude"
// @Param lon query string true "Longitude"
// @Param radius query string false "Radius in meters"
// @Success 200 {array} dto.ReportDTO
// @Failure 400 {object} util.ErrorResponse
// @Router /reports/nearby [get]
func (h *ReportHandler) ListNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	radiusStr := r.URL.Query().Get("radius")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		util.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		util.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		radius = 500 // default
	}

	list, err := h.reportUseCase.ReportUseCase.ListNearby(r.Context(), lat, lon, radius)
	if err != nil {
		slog.Error("failed to list nearby reports", "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]dto.ReportDTO, 0, len(list))
	for _, v := range list {
		response = append(response, dto.ReportToDTO(v))
	}
	util.Response(w, response, http.StatusOK)
}

// Verify godoc
// @Summary Verify a report
// @Description Verify a report by its ID
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Report ID"
// @Param verify body dto.VerifyReportRequest true "Verification data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /reports/{id}/verify [post]
func (h *ReportHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.VerifyReportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(req.ModeratorID); err != nil {
		util.Error(w, "Invalid moderator_id UUID", http.StatusBadRequest)
		return
	}

	if err := h.reportUseCase.ReportUseCase.Verify(r.Context(), id); err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{
		"status":    "verified",
		"report_id": id,
	}, http.StatusOK)
}

// Resolve godoc
// @Summary Resolve a report
// @Description Resolve a report by its ID
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Report ID"
// @Param resolve body dto.ResolveReportRequest true "Resolution data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /reports/{id}/resolve [post]
func (h *ReportHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.ResolveReportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.reportUseCase.ReportUseCase.Resolve(r.Context(), id, req.ModeratorID); err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{
		"status":    "resolved",
		"report_id": id,
	}, http.StatusOK)
}
