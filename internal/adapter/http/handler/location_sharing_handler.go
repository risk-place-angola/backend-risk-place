package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

const (
	errLocationSharingNotFound = "location sharing not found"
)

type LocationSharingHandler struct {
	app *application.Application
}

func NewLocationSharingHandler(app *application.Application) *LocationSharingHandler {
	return &LocationSharingHandler{
		app: app,
	}
}

// CreateLocationSharing godoc.
// @Summary Create location sharing session
// @Description Create a new location sharing session with expiration time (supports both authenticated and anonymous users)
// @Tags location-sharing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Device-ID header string false "Device ID for anonymous users"
// @Param request body dto.CreateLocationSharingRequest true "Location Sharing Request"
// @Success 201 {object} dto.LocationSharingResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /location-sharing [post]
func (h *LocationSharingHandler) CreateLocationSharing(w http.ResponseWriter, r *http.Request) {
	identifier, isAuthenticated := util.GetIdentifierFromContext(r.Context())
	if identifier == "" {
		slog.Error("failed to get identifier from context")
		util.Error(w, "unauthorized: JWT or X-Device-ID required", http.StatusUnauthorized)
		return
	}

	var req dto.CreateLocationSharingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var response *dto.LocationSharingResponse
	var err error

	if isAuthenticated {
		userID, parseErr := uuid.Parse(identifier)
		if parseErr != nil {
			slog.Error("invalid user ID", "error", parseErr)
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		response, err = h.app.LocationSharingUseCase.CreateLocationSharingForUser(r.Context(), userID, req)
	} else {
		response, err = h.app.LocationSharingUseCase.CreateLocationSharingForAnonymous(r.Context(), identifier, req)
	}

	if err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, response, http.StatusCreated)
}

// GetPublicLocationSharing godoc.
// @Summary Get shared location by token
// @Description Retrieve public location information using share token
// @Tags location-sharing
// @Accept json
// @Produce json
// @Param token path string true "Share Token"
// @Success 200 {object} dto.PublicLocationResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 410 {object} util.ErrorResponse
// @Router /share/{token} [get]
func (h *LocationSharingHandler) GetPublicLocationSharing(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		util.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	response, err := h.app.LocationSharingUseCase.GetLocationSharingByToken(r.Context(), token)
	if err != nil {
		if err.Error() == errLocationSharingNotFound {
			util.Error(w, errLocationSharingNotFound, http.StatusNotFound)
			return
		}
		if err.Error() == "location sharing is expired or inactive" {
			util.Error(w, "location sharing expired", http.StatusGone)
			return
		}
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, response, http.StatusOK)
}

// UpdateLocationSharing godoc.
// @Summary Update shared location coordinates
// @Description Update the current coordinates of an active location sharing session (supports both authenticated and anonymous users)
// @Tags location-sharing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Device-ID header string false "Device ID for anonymous users"
// @Param id path string true "Location Sharing ID"
// @Param request body dto.UpdateLocationRequest true "Location Update Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /location-sharing/{id}/location [put]
func (h *LocationSharingHandler) UpdateLocationSharing(w http.ResponseWriter, r *http.Request) {
	identifier, isAuthenticated := util.GetIdentifierFromContext(r.Context())
	if identifier == "" {
		slog.Error("failed to get identifier from context")
		util.Error(w, "unauthorized: JWT or X-Device-ID required", http.StatusUnauthorized)
		return
	}

	sharingIDStr := r.PathValue("id")
	if sharingIDStr == "" {
		util.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	sharingID, err := uuid.Parse(sharingIDStr)
	if err != nil {
		util.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if isAuthenticated {
		userID, parseErr := uuid.Parse(identifier)
		if parseErr != nil {
			slog.Error("invalid user ID", "error", parseErr)
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		err = h.app.LocationSharingUseCase.UpdateLocationForUser(r.Context(), sharingID, userID, req)
	} else {
		err = h.app.LocationSharingUseCase.UpdateLocationForAnonymous(r.Context(), sharingID, identifier, req)
	}

	if err != nil {
		if err.Error() == errLocationSharingNotFound {
			util.Error(w, errLocationSharingNotFound, http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized" {
			util.Error(w, "unauthorized", http.StatusForbidden)
			return
		}
		if err.Error() == "location sharing is expired or inactive" {
			util.Error(w, "location sharing expired", http.StatusGone)
			return
		}
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{"status": "location updated"}, http.StatusOK)
}

// DeleteLocationSharing godoc.
// @Summary Stop location sharing
// @Description Deactivate an active location sharing session (supports both authenticated and anonymous users)
// @Tags location-sharing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Device-ID header string false "Device ID for anonymous users"
// @Param id path string true "Location Sharing ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /location-sharing/{id} [delete]
func (h *LocationSharingHandler) DeleteLocationSharing(w http.ResponseWriter, r *http.Request) {
	identifier, isAuthenticated := util.GetIdentifierFromContext(r.Context())
	if identifier == "" {
		slog.Error("failed to get identifier from context")
		util.Error(w, "unauthorized: JWT or X-Device-ID required", http.StatusUnauthorized)
		return
	}

	sharingIDStr := r.PathValue("id")
	if sharingIDStr == "" {
		util.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	sharingID, err := uuid.Parse(sharingIDStr)
	if err != nil {
		util.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if isAuthenticated {
		userID, parseErr := uuid.Parse(identifier)
		if parseErr != nil {
			slog.Error("invalid user ID", "error", parseErr)
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		err = h.app.LocationSharingUseCase.DeleteLocationSharingForUser(r.Context(), sharingID, userID)
	} else {
		err = h.app.LocationSharingUseCase.DeleteLocationSharingForAnonymous(r.Context(), sharingID, identifier)
	}

	if err != nil {
		if err.Error() == errLocationSharingNotFound {
			util.Error(w, errLocationSharingNotFound, http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized" {
			util.Error(w, "unauthorized", http.StatusForbidden)
			return
		}
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{"status": "location sharing stopped"}, http.StatusOK)
}
