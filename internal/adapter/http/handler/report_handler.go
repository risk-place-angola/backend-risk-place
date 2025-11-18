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

// List godoc
// @Summary List all reports with pagination
// @Description List all reports in the system with pagination and filters
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param status query string false "Filter by status (pending, verified, resolved)"
// @Param sort query string false "Sort field (default: created_at)"
// @Param order query string false "Sort order (asc, desc) (default: desc)"
// @Success 200 {object} dto.ListReportsResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /reports [get]
func (h *ReportHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Convert to int with defaults
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Validate status if provided
	if status != "" && status != "pending" && status != "verified" && status != "resolved" {
		util.Error(w, "Invalid status. Must be: pending, verified, or resolved", http.StatusBadRequest)
		return
	}

	// Call use case
	params := dto.ListReportsQueryParams{
		Page:   page,
		Limit:  limit,
		Status: status,
		Sort:   sort,
		Order:  order,
	}

	response, err := h.reportUseCase.ReportUseCase.List(r.Context(), params)
	if err != nil {
		slog.Error("failed to list reports", "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, response, http.StatusOK)
}

// ListNearby godoc
// @Summary List nearby reports with distance
// @Description List reports near the specified location with calculated distance
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param latitude query number true "Latitude"
// @Param longitude query number true "Longitude"
// @Param radius query number true "Radius in meters"
// @Param limit query int false "Maximum number of results (default: 50)"
// @Success 200 {object} dto.NearbyReportsResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /reports/nearby [get]
func (h *ReportHandler) ListNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("latitude")
	lonStr := r.URL.Query().Get("longitude")
	radiusStr := r.URL.Query().Get("radius")
	limitStr := r.URL.Query().Get("limit")

	// Validate required parameters
	if latStr == "" {
		util.Error(w, "latitude is required", http.StatusBadRequest)
		return
	}
	if lonStr == "" {
		util.Error(w, "longitude is required", http.StatusBadRequest)
		return
	}
	if radiusStr == "" {
		util.Error(w, "radius is required", http.StatusBadRequest)
		return
	}

	// Parse parameters
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
		util.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	// Parse optional limit
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Call use case
	params := dto.NearbyReportsQueryParams{
		Latitude:  lat,
		Longitude: lon,
		Radius:    radius,
		Limit:     limit,
	}

	response, err := h.reportUseCase.ReportUseCase.ListNearbyWithDistance(r.Context(), params)
	if err != nil {
		slog.Error("failed to list nearby reports", "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// UpdateLocation godoc
// @Summary Update report location
// @Description Update the geographic location of a report (used when user drags marker on map)
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Report ID"
// @Param location body dto.UpdateReportLocationRequest true "New location data"
// @Success 200 {object} dto.UpdateReportLocationResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /reports/{id}/location [put]
func (h *ReportHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	reportID := r.PathValue("id")
	if reportID == "" {
		util.Error(w, "report ID is required", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(reportID); err != nil {
		slog.Error("invalid report ID format", "reportID", reportID, "error", err)
		util.Error(w, "invalid report ID format", http.StatusBadRequest)
		return
	}

	var req dto.UpdateReportLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if h.reportUseCase == nil || h.reportUseCase.ReportUseCase == nil {
		slog.Error("reportUseCase is nil")
		util.Error(w, "internal server error: use case not initialized", http.StatusInternalServerError)
		return
	}

	if err := h.reportUseCase.ReportUseCase.UpdateLocation(r.Context(), reportID, req); err != nil {
		slog.Error("failed to update report location", "reportID", reportID, "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{
		"id":      reportID,
		"message": "Report location updated successfully",
	}, http.StatusOK)
}
