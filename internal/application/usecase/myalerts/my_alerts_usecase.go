package myalerts

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type MyAlertsUseCase struct {
	alertRepo     repository.AlertRepository
	riskTypeRepo  repository.RiskTypesRepository
	riskTopicRepo repository.RiskTopicsRepository
}

func NewMyAlertsUseCase(
	alertRepo repository.AlertRepository,
	riskTypeRepo repository.RiskTypesRepository,
	riskTopicRepo repository.RiskTopicsRepository,
) *MyAlertsUseCase {
	return &MyAlertsUseCase{
		alertRepo:     alertRepo,
		riskTypeRepo:  riskTypeRepo,
		riskTopicRepo: riskTopicRepo,
	}
}

func (uc *MyAlertsUseCase) GetMyCreatedAlerts(ctx context.Context, userID uuid.UUID) ([]dto.MyAlertResponse, error) {
	alerts, err := uc.alertRepo.GetByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching user alerts", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch alerts")
	}

	return uc.toResponseList(ctx, alerts, userID)
}

func (uc *MyAlertsUseCase) GetMySubscribedAlerts(ctx context.Context, userID uuid.UUID) ([]dto.MyAlertResponse, error) {
	alerts, err := uc.alertRepo.GetSubscribedAlerts(ctx, userID)
	if err != nil {
		slog.Error("Error fetching subscribed alerts", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch subscribed alerts")
	}

	return uc.toResponseList(ctx, alerts, userID)
}

func (uc *MyAlertsUseCase) UpdateAlert(ctx context.Context, userID, alertID uuid.UUID, input dto.UpdateAlertInput) (*dto.MyAlertResponse, error) {
	alert, err := uc.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		slog.Error("Error fetching alert for update", "alert_id", alertID, "error", err)
		return nil, errors.New("failed to fetch alert")
	}

	if alert == nil {
		return nil, errors.New("alert not found")
	}

	if alert.CreatedBy == nil || *alert.CreatedBy != userID {
		return nil, errors.New("unauthorized: you can only update your own alerts")
	}

	alert.Message = input.Message
	alert.Severity = model.Severity(input.Severity)
	alert.RadiusMeters = input.RadiusMeters

	if err := uc.alertRepo.Update(ctx, alert); err != nil {
		slog.Error("Error updating alert", "alert_id", alertID, "error", err)
		return nil, errors.New("failed to update alert")
	}

	return uc.toResponse(ctx, alert, userID)
}

func (uc *MyAlertsUseCase) DeleteAlert(ctx context.Context, userID, alertID uuid.UUID) error {
	alert, err := uc.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		slog.Error("Error fetching alert for deletion", "alert_id", alertID, "error", err)
		return errors.New("failed to fetch alert")
	}

	if alert == nil {
		return errors.New("alert not found")
	}

	if alert.CreatedBy == nil || *alert.CreatedBy != userID {
		return errors.New("unauthorized: you can only delete your own alerts")
	}

	if err := uc.alertRepo.Delete(ctx, alertID, userID); err != nil {
		slog.Error("Error deleting alert", "alert_id", alertID, "error", err)
		return errors.New("failed to delete alert")
	}

	return nil
}

func (uc *MyAlertsUseCase) SubscribeToAlert(ctx context.Context, userID, alertID uuid.UUID) (*dto.AlertSubscriptionResponse, error) {
	// Check if already subscribed
	isSubscribed, err := uc.alertRepo.IsUserSubscribed(ctx, alertID, userID)
	if err != nil {
		slog.Error("Error checking subscription", "alert_id", alertID, "user_id", userID, "error", err)
		return nil, errors.New("failed to check subscription")
	}

	if isSubscribed {
		return &dto.AlertSubscriptionResponse{
			Success: true,
			Message: "Already subscribed to alert",
		}, nil
	}

	alert, err := uc.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		slog.Error("Error fetching alert for subscription", "alert_id", alertID, "error", err)
		return nil, errors.New("failed to fetch alert")
	}

	if alert == nil {
		return nil, errors.New("alert not found")
	}

	if alert.CreatedBy != nil && *alert.CreatedBy == userID {
		return nil, errors.New("you cannot subscribe to your own alert")
	}

	subscription, err := model.NewAlertSubscription(alertID, userID)
	if err != nil {
		return nil, err
	}

	if err := uc.alertRepo.SubscribeToAlert(ctx, subscription); err != nil {
		slog.Error("Error subscribing to alert", "alert_id", alertID, "user_id", userID, "error", err)
		return nil, errors.New("failed to subscribe to alert")
	}

	return &dto.AlertSubscriptionResponse{
		Success: true,
		Message: "Successfully subscribed to alert",
	}, nil
}

func (uc *MyAlertsUseCase) UnsubscribeFromAlert(ctx context.Context, userID, alertID uuid.UUID) (*dto.AlertSubscriptionResponse, error) {
	isSubscribed, err := uc.alertRepo.IsUserSubscribed(ctx, alertID, userID)
	if err != nil {
		slog.Error("Error checking subscription", "alert_id", alertID, "user_id", userID, "error", err)
		return nil, errors.New("failed to check subscription")
	}

	if !isSubscribed {
		return nil, errors.New("you are not subscribed to this alert")
	}

	if err := uc.alertRepo.UnsubscribeFromAlert(ctx, alertID, userID); err != nil {
		slog.Error("Error unsubscribing from alert", "alert_id", alertID, "user_id", userID, "error", err)
		return nil, errors.New("failed to unsubscribe from alert")
	}

	return &dto.AlertSubscriptionResponse{
		Success: true,
		Message: "Successfully unsubscribed from alert",
	}, nil
}

func (uc *MyAlertsUseCase) toResponseList(ctx context.Context, alerts []*model.Alert, userID uuid.UUID) ([]dto.MyAlertResponse, error) {
	response := make([]dto.MyAlertResponse, 0, len(alerts))

	for _, alert := range alerts {
		alertResponse, err := uc.toResponse(ctx, alert, userID)
		if err != nil {
			slog.Error("Error converting alert to response", "alert_id", alert.ID, "error", err)
			continue
		}
		response = append(response, *alertResponse)
	}

	return response, nil
}

func (uc *MyAlertsUseCase) toResponse(ctx context.Context, alert *model.Alert, userID uuid.UUID) (*dto.MyAlertResponse, error) {
	response := &dto.MyAlertResponse{
		ID:           alert.ID.String(),
		Message:      alert.Message,
		Latitude:     alert.Latitude,
		Longitude:    alert.Longitude,
		Province:     alert.Province,
		Municipality: alert.Municipality,
		Neighborhood: alert.Neighborhood,
		Address:      alert.Address,
		RadiusMeters: alert.RadiusMeters,
		Status:       string(alert.Status),
		Severity:     string(alert.Severity),
		CreatedAt:    alert.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		RiskTypeName: alert.RiskTypeName,
	}

	// Convert icon paths to URLs
	if alert.RiskTypeIconPath != nil && *alert.RiskTypeIconPath != "" {
		iconURL := "/api/v1/storage/" + *alert.RiskTypeIconPath
		response.RiskTypeIconURL = &iconURL
	}

	if alert.RiskTopicID != uuid.Nil {
		response.RiskTopicName = alert.RiskTopicName

		// Convert topic icon path to URL
		if alert.RiskTopicIconPath != nil && *alert.RiskTopicIconPath != "" {
			topicIconURL := "/api/v1/storage/" + *alert.RiskTopicIconPath
			response.RiskTopicIconURL = &topicIconURL
		}
	}

	if !alert.ExpiresAt.IsZero() {
		response.ExpiresAt = alert.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if !alert.ResolvedAt.IsZero() {
		response.ResolvedAt = alert.ResolvedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	isSubscribed, err := uc.alertRepo.IsUserSubscribed(ctx, alert.ID, userID)
	if err == nil {
		response.IsSubscribed = isSubscribed
	}

	return response, nil
}
