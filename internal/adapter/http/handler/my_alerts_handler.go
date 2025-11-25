package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type MyAlertsHandler struct {
	app                  *application.Application
	anonymousSessionRepo repository.AnonymousSessionRepository
	queries              *sqlc.Queries
}

func NewMyAlertsHandler(app *application.Application, anonymousSessionRepo repository.AnonymousSessionRepository, queries *sqlc.Queries) *MyAlertsHandler {
	return &MyAlertsHandler{
		app:                  app,
		anonymousSessionRepo: anonymousSessionRepo,
		queries:              queries,
	}
}

// GetMyCreatedAlerts godoc
// @Summary Get all alerts created by the current user
// @Description Retrieve all alerts that were created by the authenticated user
// @Tags my-alerts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.MyAlertResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/alerts/created [get]
func (h *MyAlertsHandler) GetMyCreatedAlerts(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := dto.ParseUUID(userIDStr)
	if err != nil {
		slog.Error("invalid user ID in context", "error", err)
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	alerts, err := h.app.MyAlertsUseCase.GetMyCreatedAlerts(r.Context(), uid)
	if err != nil {
		slog.Error("error fetching user alerts", "user_id", uid, "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, alerts, http.StatusOK)
}

// GetMySubscribedAlerts godoc
// @Summary Get all alerts the user is subscribed to
// @Description Retrieve all alerts that the authenticated user has subscribed to
// @Tags my-alerts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.MyAlertResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/alerts/subscribed [get]
func (h *MyAlertsHandler) GetMySubscribedAlerts(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := dto.ParseUUID(userIDStr)
	if err != nil {
		slog.Error("invalid user ID in context", "error", err)
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	alerts, err := h.app.MyAlertsUseCase.GetMySubscribedAlerts(r.Context(), uid)
	if err != nil {
		slog.Error("error fetching subscribed alerts", "user_id", uid, "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, alerts, http.StatusOK)
}

// UpdateAlert godoc
// @Summary Update an alert
// @Description Update an alert created by the authenticated user
// @Tags my-alerts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Param alert body dto.UpdateAlertInput true "Updated alert data"
// @Success 200 {object} dto.MyAlertResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /alerts/{id} [put]
func (h *MyAlertsHandler) UpdateAlert(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := dto.ParseUUID(userIDStr)
	if err != nil {
		slog.Error("invalid user ID in context", "error", err)
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	alertID := r.PathValue("id")
	if alertID == "" {
		util.Error(w, "alert ID is required", http.StatusBadRequest)
		return
	}

	aid, err := uuid.Parse(alertID)
	if err != nil {
		util.Error(w, "invalid alert ID", http.StatusBadRequest)
		return
	}

	var input dto.UpdateAlertInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	alert, err := h.app.MyAlertsUseCase.UpdateAlert(r.Context(), uid, aid, input)
	if err != nil {
		slog.Error("error updating alert", "user_id", uid, "alert_id", aid, "error", err)
		if errors.Is(err, domainErrors.ErrAlertNotFound) {
			util.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: you can only update your own alerts" {
			util.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.Response(w, alert, http.StatusOK)
}

// DeleteAlert godoc
// @Summary Delete an alert
// @Description Delete an alert created by the authenticated user
// @Tags my-alerts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Alert ID"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /alerts/{id} [delete]
func (h *MyAlertsHandler) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	alertID := r.PathValue("id")
	if alertID == "" {
		util.Error(w, "alert ID is required", http.StatusBadRequest)
		return
	}

	aid, err := uuid.Parse(alertID)
	if err != nil {
		util.Error(w, "invalid alert ID", http.StatusBadRequest)
		return
	}

	if err := h.app.MyAlertsUseCase.DeleteAlert(r.Context(), uid, aid); err != nil {
		slog.Error("error deleting alert", "user_id", uid, "alert_id", aid, "error", err)
		if errors.Is(err, domainErrors.ErrAlertNotFound) {
			util.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: you can only delete your own alerts" {
			util.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SubscribeToAlert godoc
// @Summary Subscribe to an alert
// @Description Subscribe to receive notifications for an alert
// @Tags my-alerts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} dto.AlertSubscriptionResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /alerts/{id}/subscribe [post]
func (h *MyAlertsHandler) SubscribeToAlert(w http.ResponseWriter, r *http.Request) {
	identifier, ok := util.ExtractUserIdentifierOrError(w, r)
	if !ok {
		return
	}

	aid, ok := util.ExtractAndValidatePathID(w, r, "id", "alert")
	if !ok {
		return
	}

	if identifier.IsAuthenticated {
		h.subscribeAuthenticatedUser(w, r, aid, identifier.UserID)
	} else {
		h.subscribeAnonymousUser(w, r, aid, identifier.DeviceID)
	}
}

// UnsubscribeFromAlert godoc
// @Summary Unsubscribe from an alert
// @Description Unsubscribe from receiving notifications for an alert
// @Tags my-alerts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} dto.AlertSubscriptionResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /alerts/{id}/unsubscribe [delete]
func (h *MyAlertsHandler) UnsubscribeFromAlert(w http.ResponseWriter, r *http.Request) {
	identifier, ok := util.ExtractUserIdentifierOrError(w, r)
	if !ok {
		return
	}

	aid, ok := util.ExtractAndValidatePathID(w, r, "id", "alert")
	if !ok {
		return
	}

	if identifier.IsAuthenticated {
		h.unsubscribeAuthenticatedUser(w, r, aid, identifier.UserID)
	} else {
		h.unsubscribeAnonymousUser(w, r, aid, identifier.DeviceID)
	}
}

func (h *MyAlertsHandler) subscribeAuthenticatedUser(w http.ResponseWriter, r *http.Request, aid uuid.UUID, userIDStr string) {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := h.app.MyAlertsUseCase.SubscribeToAlert(r.Context(), uid, aid)
	if err != nil {
		slog.Error("failed to subscribe to alert", "user_id", uid, "alert_id", aid, "error", err)
		if errors.Is(err, domainErrors.ErrAlertNotFound) {
			util.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.Response(w, result, http.StatusOK)
}

func (h *MyAlertsHandler) subscribeAnonymousUser(w http.ResponseWriter, r *http.Request, aid uuid.UUID, deviceID string) {
	subscribed, err := h.queries.IsAnonymousSubscribedToAlert(r.Context(), sqlc.IsAnonymousSubscribedToAlertParams{
		AlertID:  aid,
		DeviceID: sql.NullString{String: deviceID, Valid: true},
	})
	if err != nil {
		slog.Error("failed to check subscription status", "device_id", deviceID, "alert_id", aid, "error", err)
		util.Error(w, "failed to check subscription", http.StatusInternalServerError)
		return
	}

	if subscribed {
		util.Response(w, map[string]interface{}{
			"message":    "already subscribed",
			"alert_id":   aid,
			"device_id":  deviceID,
			"subscribed": true,
		}, http.StatusOK)
		return
	}

	session, err := h.getOrCreateSession(r, deviceID)
	if err != nil {
		util.Error(w, "failed to get or create session", http.StatusInternalServerError)
		return
	}

	err = h.queries.SubscribeAnonymousToAlert(r.Context(), sqlc.SubscribeAnonymousToAlertParams{
		ID:                 uuid.New(),
		AlertID:            aid,
		AnonymousSessionID: uuid.NullUUID{UUID: session.ID, Valid: true},
		DeviceID:           sql.NullString{String: deviceID, Valid: true},
		SubscribedAt:       time.Now(),
	})

	if err != nil {
		slog.Error("failed to subscribe anonymous to alert", "device_id", deviceID, "alert_id", aid, "error", err)
		util.Error(w, "failed to subscribe to alert", http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]interface{}{
		"success": true,
		"message": "Successfully subscribed to alert",
	}, http.StatusOK)
}

func (h *MyAlertsHandler) getOrCreateSession(r *http.Request, deviceID string) (*model.AnonymousSession, error) {
	session, err := h.anonymousSessionRepo.FindByDeviceID(r.Context(), deviceID)
	if err != nil {
		session, err = model.NewAnonymousSession(deviceID, "", "", "")
		if err != nil {
			slog.Error("failed to create anonymous session", "device_id", deviceID, "error", err)
			return nil, err
		}

		if err := h.anonymousSessionRepo.Create(r.Context(), session); err != nil {
			slog.Error("failed to save anonymous session", "device_id", deviceID, "error", err)
			return nil, err
		}
	}
	return session, nil
}

func (h *MyAlertsHandler) unsubscribeAuthenticatedUser(w http.ResponseWriter, r *http.Request, aid uuid.UUID, userIDStr string) {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := h.app.MyAlertsUseCase.UnsubscribeFromAlert(r.Context(), uid, aid)
	if err != nil {
		slog.Error("error unsubscribing from alert", "user_id", uid, "alert_id", aid, "error", err)
		if err.Error() == "you are not subscribed to this alert" {
			util.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, result, http.StatusOK)
}

func (h *MyAlertsHandler) unsubscribeAnonymousUser(w http.ResponseWriter, r *http.Request, aid uuid.UUID, deviceID string) {
	session, err := h.anonymousSessionRepo.FindByDeviceID(r.Context(), deviceID)
	if err != nil {
		slog.Error("anonymous session not found", "device_id", deviceID, "error", err)
		util.Error(w, "session not found", http.StatusNotFound)
		return
	}

	err = h.queries.UnsubscribeAnonymousFromAlert(r.Context(), sqlc.UnsubscribeAnonymousFromAlertParams{
		AlertID:            aid,
		AnonymousSessionID: uuid.NullUUID{UUID: session.ID, Valid: true},
		DeviceID:           sql.NullString{String: deviceID, Valid: true},
	})

	if err != nil {
		slog.Error("failed to unsubscribe anonymous from alert", "device_id", deviceID, "alert_id", aid, "error", err)
		util.Error(w, "failed to unsubscribe from alert", http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]interface{}{
		"success": true,
		"message": "Successfully unsubscribed from alert",
	}, http.StatusOK)
}
