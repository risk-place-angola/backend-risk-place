package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type AlertHandler struct {
	alertUseCase         *application.Application
	anonymousSessionRepo repository.AnonymousSessionRepository
	queries              *sqlc.Queries
}

func NewAlertHandler(alertUseCase *application.Application, anonymousSessionRepo repository.AnonymousSessionRepository, queries *sqlc.Queries) *AlertHandler {
	return &AlertHandler{
		alertUseCase:         alertUseCase,
		anonymousSessionRepo: anonymousSessionRepo,
		queries:              queries,
	}
}

// CreateAlert godoc.
// @Summary Create a new alert.
// @Description Create a new alert (supports both authenticated and anonymous users).
// @Tags alerts
// @Accept json
// @Produce json
// @Security OptionalAuth
// @Param X-Device-Id header string false "Device ID for anonymous users"
// @Param alert body dto.Alert true "Alert"
// @Success 201 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /alerts [post].
func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	identifier, ok := util.ExtractUserIdentifierOrError(w, r)
	if !ok {
		return
	}

	var req dto.Alert
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if identifier.IsAuthenticated {
		h.createAuthenticatedAlert(w, r, req)
	} else {
		h.createAnonymousAlert(w, r, req, identifier.DeviceID)
	}
}

func (h *AlertHandler) createAuthenticatedAlert(w http.ResponseWriter, r *http.Request, req dto.Alert) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	req.UserID = userIDStr
	err := h.alertUseCase.AlertUseCase.TriggerAlert(r.Context(), req)
	if err != nil {
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	util.Response(w, map[string]string{"status": "alert triggered"}, http.StatusCreated)
}

func (h *AlertHandler) createAnonymousAlert(w http.ResponseWriter, r *http.Request, req dto.Alert, deviceID string) {
	session, err := h.getOrCreateAnonymousSession(r, deviceID)
	if err != nil {
		util.Error(w, "failed to get or create session", http.StatusInternalServerError)
		return
	}

	riskTypeID, parseErr := uuid.Parse(req.RiskTypeID)
	if parseErr != nil {
		util.Error(w, "invalid risk type ID", http.StatusBadRequest)
		return
	}

	var riskTopicIDNullUUID uuid.NullUUID
	if req.RiskTopicID != "" {
		topicID, topicErr := uuid.Parse(req.RiskTopicID)
		if topicErr != nil {
			util.Error(w, "invalid risk topic ID", http.StatusBadRequest)
			return
		}
		riskTopicIDNullUUID = uuid.NullUUID{UUID: topicID, Valid: true}
	}

	err = h.queries.CreateAnonymousAlert(r.Context(), sqlc.CreateAnonymousAlertParams{
		ID:                 uuid.New(),
		AnonymousSessionID: uuid.NullUUID{UUID: session.ID, Valid: true},
		DeviceID:           sql.NullString{String: deviceID, Valid: true},
		RiskTypeID:         riskTypeID,
		RiskTopicID:        riskTopicIDNullUUID,
		Message:            req.Message,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		Province:           sql.NullString{},
		Municipality:       sql.NullString{},
		Neighborhood:       sql.NullString{},
		Address:            sql.NullString{},
		RadiusMeters:       int32(req.Radius),
		Severity:           req.Severity,
		ExpiresAt:          sql.NullTime{},
	})
	if err != nil {
		slog.Error("failed to create anonymous alert", "error", err)
		util.Error(w, "failed to create alert", http.StatusInternalServerError)
		return
	}
	util.Response(w, map[string]string{"status": "alert triggered"}, http.StatusCreated)
}

func (h *AlertHandler) getOrCreateAnonymousSession(r *http.Request, deviceID string) (*model.AnonymousSession, error) {
	session, err := h.anonymousSessionRepo.FindByDeviceID(r.Context(), deviceID)
	if err != nil {
		newSession, createErr := model.NewAnonymousSession(deviceID, "", "", "")
		if createErr != nil {
			slog.Error("failed to create anonymous session", "error", createErr)
			return nil, createErr
		}

		if createErr := h.anonymousSessionRepo.Create(r.Context(), newSession); createErr != nil {
			slog.Error("failed to save anonymous session", "error", createErr)
			return nil, createErr
		}

		return newSession, nil
	}
	return session, nil
}
