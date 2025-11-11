package alert

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/event"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type AlertUseCase struct {
	locationStore   port.LocationStore
	geoService      port.GeolocationService
	repo            repository.AlertRepository
	riskTypesRepo   repository.RiskTypesRepository
	eventDispatcher port.EventDispatcher
}

func NewAlertUseCase(
	locationStore port.LocationStore,
	geoService port.GeolocationService,
	repo repository.AlertRepository,
	riskTypesRepo repository.RiskTypesRepository,
	eventDispatcher port.EventDispatcher,
) *AlertUseCase {
	return &AlertUseCase{
		locationStore:   locationStore,
		geoService:      geoService,
		repo:            repo,
		riskTypesRepo:   riskTypesRepo,
		eventDispatcher: eventDispatcher,
	}
}

// TriggerAlert sends an alert notification to all users within the specified radius.
func (uc *AlertUseCase) TriggerAlert(ctx context.Context, alert dto.Alert) error {
	alrt := &model.Alert{
		ID:           uuid.New(),
		RiskTypeID:   uuid.MustParse(alert.RiskTypeID),
		RiskTopicID:  uuid.MustParse(alert.RiskTopicID),
		Message:      alert.Message,
		Latitude:     alert.Latitude,
		Longitude:    alert.Longitude,
		CreatedBy:    uuid.MustParse(alert.UserID),
		RadiusMeters: int(alert.Radius),
		Severity:     model.Severity(alert.Severity),
		Status:       model.AlertStatusActive,
		CreatedAt:    time.Now(),
	}

	err := uc.geoService.ValidateCoordinates(alrt.Latitude, alrt.Longitude)
	if err != nil {
		slog.Error("invalid coordinates for alert", "error", err)
		return err
	}

	if alert.Radius <= 0 {
		riskType, err := uc.riskTypesRepo.GetRiskTypeByID(ctx, alert.RiskTypeID)
		if err != nil {
			slog.Error("failed to get risk type for alert", "error", err)
			return err
		}
		alrt.RadiusMeters = riskType.DefaultRadiusMeters
	}

	err = uc.repo.Create(ctx, alrt)
	if err != nil {
		slog.Error("failed to create alert", "error", err)
		return err
	}

	userIDs, err := uc.locationStore.FindUsersInRadius(ctx, alert.Latitude, alert.Longitude, alert.Radius)
	if err != nil {
		slog.Error("failed to find users in radius", "error", err)
		return err
	}

	uuidUserIDs := make([]uuid.UUID, 0, len(userIDs))
	for _, uid := range userIDs {
		if err := uc.repo.CreateAlertNotification(ctx, alrt.ID, uid); err != nil {
			slog.Error("failed to create alert notification", "error", err, "user_id", uid)
			continue
		}
		slog.Info("created alert notification", "alert_id", alrt.ID.String(), "user_id", uid)
		uuidUserIDs = append(uuidUserIDs, uuid.MustParse(uid))
	}

	uc.eventDispatcher.Dispatch(event.AlertCreatedEvent{
		AlertID:   alrt.ID,
		UserID:    uuidUserIDs,
		Message:   alert.Message,
		Latitude:  alert.Latitude,
		Longitude: alert.Longitude,
		Radius:    alert.Radius,
	})

	return nil
}
