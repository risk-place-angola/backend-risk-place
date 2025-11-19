package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

type SafeRouteHandler struct {
	app *application.Application
}

func NewSafeRouteHandler(app *application.Application) *SafeRouteHandler {
	return &SafeRouteHandler{
		app: app,
	}
}

// decodeAndHandle is a helper to reduce code duplication in handlers
func (h *SafeRouteHandler) decodeAndHandle(
	w http.ResponseWriter,
	r *http.Request,
	req interface{},
	handler func() (interface{}, error),
	errorMsg string,
) {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := handler()
	if err != nil {
		slog.Error(errorMsg, "error", err)
		util.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	util.Response(w, response, http.StatusOK)
}

func (h *SafeRouteHandler) CalculateSafeRoute(w http.ResponseWriter, r *http.Request) {
	var req dto.SafeRouteRequest
	h.decodeAndHandle(w, r, &req, func() (interface{}, error) {
		return h.app.SafeRouteUseCase.CalculateSafeRoute(r.Context(), &req)
	}, "failed to calculate safe route")
}

func (h *SafeRouteHandler) GetIncidentsHeatmap(w http.ResponseWriter, r *http.Request) {
	var req dto.HeatmapRequest
	h.decodeAndHandle(w, r, &req, func() (interface{}, error) {
		return h.app.SafeRouteUseCase.GetIncidentsHeatmap(r.Context(), &req)
	}, "failed to get incidents heatmap")
}

// navigateToSavedLocation is a helper to reduce duplication between NavigateToHome and NavigateToWork
func (h *SafeRouteHandler) navigateToSavedLocation(
	w http.ResponseWriter,
	r *http.Request,
	useCase func(r *http.Request, currentLat, currentLon float64) (*dto.SafeRouteResponse, error),
	notConfiguredMsg, errorMsg string,
) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := dto.ParseUUID(userIDStr)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusUnauthorized)
		return
	}

	var req dto.NavigateToSavedLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := useCase(r, req.CurrentLat, req.CurrentLon)
	if err != nil {
		if err.Error() == notConfiguredMsg {
			util.Error(w, notConfiguredMsg, http.StatusNotFound)
			return
		}
		slog.Error(errorMsg, "error", err)
		util.Error(w, "failed to calculate route", http.StatusInternalServerError)
		return
	}

	util.Response(w, response, http.StatusOK)
}

// NavigateToHome godoc
// @Summary Calculate safe route to home address
// @Description Calculate a safe route from current location to the user's saved home address
// @Tags routes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param location body dto.NavigateToSavedLocationRequest true "Current location coordinates"
// @Success 200 {object} dto.SafeRouteResponse "Safe route calculated successfully"
// @Failure 400 {object} util.ErrorResponse "Invalid request body"
// @Failure 401 {object} util.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} util.ErrorResponse "Home address not configured"
// @Failure 500 {object} util.ErrorResponse "Failed to calculate route"
// @Router /routes/navigate-home [post]
func (h *SafeRouteHandler) NavigateToHome(w http.ResponseWriter, r *http.Request) {
	h.navigateToSavedLocation(w, r,
		func(req *http.Request, currentLat, currentLon float64) (*dto.SafeRouteResponse, error) {
			userIDStr, _ := util.GetUserIDFromContext(req.Context())
			userID, _ := dto.ParseUUID(userIDStr)
			return h.app.SafeRouteUseCase.NavigateToHome(req.Context(), userID, currentLat, currentLon)
		},
		"home address not configured",
		"failed to navigate to home",
	)
}

// NavigateToWork godoc
// @Summary Calculate safe route to work address
// @Description Calculate a safe route from current location to the user's saved work address
// @Tags routes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param location body dto.NavigateToSavedLocationRequest true "Current location coordinates"
// @Success 200 {object} dto.SafeRouteResponse "Safe route calculated successfully"
// @Failure 400 {object} util.ErrorResponse "Invalid request body"
// @Failure 401 {object} util.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} util.ErrorResponse "Work address not configured"
// @Failure 500 {object} util.ErrorResponse "Failed to calculate route"
// @Router /routes/navigate-work [post]
func (h *SafeRouteHandler) NavigateToWork(w http.ResponseWriter, r *http.Request) {
	h.navigateToSavedLocation(w, r,
		func(req *http.Request, currentLat, currentLon float64) (*dto.SafeRouteResponse, error) {
			userIDStr, _ := util.GetUserIDFromContext(req.Context())
			userID, _ := dto.ParseUUID(userIDStr)
			return h.app.SafeRouteUseCase.NavigateToWork(req.Context(), userID, currentLat, currentLon)
		},
		"work address not configured",
		"failed to navigate to work",
	)
}
