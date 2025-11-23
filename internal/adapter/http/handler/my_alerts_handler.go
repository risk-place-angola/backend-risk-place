package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
)

type MyAlertsHandler struct {
	app *application.Application
}

func NewMyAlertsHandler(app *application.Application) *MyAlertsHandler {
	return &MyAlertsHandler{
		app: app,
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
	uid, ok := util.ExtractAndValidateUserID(w, r)
	if !ok {
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
	uid, ok := util.ExtractAndValidateUserID(w, r)
	if !ok {
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
	uid, ok := util.ExtractAndValidateUserID(w, r)
	if !ok {
		return
	}

	aid, ok := util.ExtractAndValidatePathID(w, r, "id", "alert")
	if !ok {
		return
	}

	result, err := h.app.MyAlertsUseCase.SubscribeToAlert(r.Context(), uid, aid)
	if err != nil {
		slog.Error("error subscribing to alert", "user_id", uid, "alert_id", aid, "error", err)
		if errors.Is(err, domainErrors.ErrAlertNotFound) {
			util.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.Response(w, result, http.StatusOK)
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
	uid, ok := util.ExtractAndValidateUserID(w, r)
	if !ok {
		return
	}

	aid, ok := util.ExtractAndValidatePathID(w, r, "id", "alert")
	if !ok {
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
