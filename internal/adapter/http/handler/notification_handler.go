package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

type NotificationHandler struct {
	app *application.Application
}

func NewNotificationHandler(app *application.Application) *NotificationHandler {
	return &NotificationHandler{
		app: app,
	}
}

// UpdateDeviceInfo godoc
// @Summary Update user device information
// @Description Update FCM token and device language for push notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param device body dto.UpdateDeviceInfoRequest true "Device information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/device [put]
func (h *NotificationHandler) UpdateDeviceInfo(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateDeviceInfoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := dto.ParseUUID(userIDStr)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.app.UserUseCase.UpdateDeviceInfo(r.Context(), uid, req.DeviceFCMToken, req.DeviceLanguage)
	if err != nil {
		util.Error(w, "failed to update device info", http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{"message": "device info updated successfully"}, http.StatusOK)
}

// UpdateNotificationPreferences godoc
// @Summary Update notification preferences
// @Description Update push and SMS notification preferences for authenticated users or anonymous sessions
// @Tags notifications
// @Accept json
// @Produce json
// @Param X-Device-Id header string false "Device ID for anonymous users"
// @Param preferences body dto.NotificationPreferencesRequest true "Notification preferences"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/notifications/preferences [put]
func (h *NotificationHandler) UpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	var req dto.NotificationPreferencesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userIDStr, hasUser := util.GetUserIDFromContext(r.Context())
	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = r.Header.Get("Device-Id")
	}

	switch {
	case hasUser:
		uid, err := dto.ParseUUID(userIDStr)
		if err != nil {
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		err = h.app.UserUseCase.UpdateNotificationPreferences(r.Context(), uid, "", req.PushEnabled, req.SMSEnabled)
		if err != nil {
			util.Error(w, "failed to update preferences", http.StatusInternalServerError)
			return
		}
	case deviceID != "":
		err := h.app.UserUseCase.UpdateNotificationPreferences(r.Context(), uuid.Nil, deviceID, req.PushEnabled, req.SMSEnabled)
		if err != nil {
			util.Error(w, "failed to update preferences", http.StatusInternalServerError)
			return
		}
	default:
		util.Error(w, "user ID or device ID required", http.StatusBadRequest)
		return
	}

	util.Response(w, map[string]string{"message": "notification preferences updated successfully"}, http.StatusOK)
}

// GetNotificationPreferences godoc
// @Summary Get notification preferences
// @Description Get push and SMS notification preferences for authenticated users or anonymous sessions
// @Tags notifications
// @Produce json
// @Param X-Device-Id header string false "Device ID for anonymous users"
// @Success 200 {object} dto.NotificationPreferencesResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/notifications/preferences [get]
func (h *NotificationHandler) GetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	userIDStr, hasUser := util.GetUserIDFromContext(r.Context())
	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = r.Header.Get("Device-Id")
	}

	var pushEnabled, smsEnabled bool
	var err error

	switch {
	case hasUser:
		uid, parseErr := dto.ParseUUID(userIDStr)
		if parseErr != nil {
			util.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		pushEnabled, smsEnabled, err = h.app.UserUseCase.GetNotificationPreferences(r.Context(), uid, "")
		if err != nil {
			util.Error(w, "failed to get preferences", http.StatusInternalServerError)
			return
		}
	case deviceID != "":
		pushEnabled, smsEnabled, err = h.app.UserUseCase.GetNotificationPreferences(r.Context(), uuid.Nil, deviceID)
		if err != nil {
			util.Error(w, "failed to get preferences", http.StatusInternalServerError)
			return
		}
	default:
		util.Error(w, "user ID or device ID required", http.StatusBadRequest)
		return
	}

	util.Response(w, dto.NotificationPreferencesResponse{
		PushEnabled: pushEnabled,
		SMSEnabled:  smsEnabled,
	}, http.StatusOK)
}
