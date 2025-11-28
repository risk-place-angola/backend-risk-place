package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type SafetySettingsHandler struct {
	app                  *application.Application
	anonymousSessionRepo repository.AnonymousSessionRepository
}

func NewSafetySettingsHandler(app *application.Application, anonymousSessionRepo repository.AnonymousSessionRepository) *SafetySettingsHandler {
	return &SafetySettingsHandler{
		app:                  app,
		anonymousSessionRepo: anonymousSessionRepo,
	}
}

// GetSettings godoc
// @Summary Get user safety settings
// @Description Retrieve safety settings for the authenticated user or anonymous user. Creates default settings if none exist.
// @Tags safety-settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.SafetySettingsResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/settings [get]
func (h *SafetySettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	identifier, ok := util.ExtractUserIdentifierOrError(w, r)
	if !ok {
		return
	}

	var settings *dto.SafetySettingsResponse
	var err error

	if identifier.IsAuthenticated {
		uid, parseErr := dto.ParseUUID(identifier.UserID)
		if parseErr != nil {
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		settings, err = h.app.SafetySettingsUseCase.GetSettings(r.Context(), uid)
	} else {
		settings, err = h.app.SafetySettingsUseCase.GetSettingsByDeviceID(r.Context(), identifier.DeviceID, h.anonymousSessionRepo)
	}

	if err != nil {
		slog.Error("error fetching safety settings", "identifier", identifier, "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, settings, http.StatusOK)
}

// UpdateSettings godoc
// @Summary Update user safety settings
// @Description Update safety settings for the authenticated user or anonymous user. All fields are optional.
// @Tags safety-settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param settings body dto.UpdateSafetySettingsInput true "Updated settings"
// @Success 200 {object} dto.SafetySettingsResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/settings [put]
func (h *SafetySettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	identifier, ok := util.ExtractUserIdentifierOrError(w, r)
	if !ok {
		return
	}

	var input dto.UpdateSafetySettingsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var settings *dto.SafetySettingsResponse
	var err error

	if identifier.IsAuthenticated {
		uid, parseErr := dto.ParseUUID(identifier.UserID)
		if parseErr != nil {
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		settings, err = h.app.SafetySettingsUseCase.UpdateSettings(r.Context(), uid, input)
	} else {
		settings, err = h.app.SafetySettingsUseCase.UpdateSettingsByDeviceID(r.Context(), identifier.DeviceID, input, h.anonymousSessionRepo)
	}
	if err != nil {
		slog.Error("error updating safety settings", "identifier", identifier, "error", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid profile_visibility value" ||
			err.Error() == "notification_alert_radius_mins must be between 100 and 10000" ||
			err.Error() == "notification_report_radius_mins must be between 100 and 10000" ||
			err.Error() == "invalid high_risk_start_time format, expected HH:MM" ||
			err.Error() == "invalid high_risk_end_time format, expected HH:MM" ||
			err.Error() == "invalid night_mode_start_time format, expected HH:MM" ||
			err.Error() == "invalid night_mode_end_time format, expected HH:MM" {
			statusCode = http.StatusBadRequest
		}
		util.Error(w, err.Error(), statusCode)
		return
	}

	util.Response(w, settings, http.StatusOK)
}
