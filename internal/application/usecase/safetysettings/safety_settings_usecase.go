package safetysettings

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type SafetySettingsUseCase struct {
	repo repository.SafetySettingsRepository
}

func NewSafetySettingsUseCase(repo repository.SafetySettingsRepository) *SafetySettingsUseCase {
	return &SafetySettingsUseCase{repo: repo}
}

func (uc *SafetySettingsUseCase) GetSettings(ctx context.Context, userID uuid.UUID) (*dto.SafetySettingsResponse, error) {
	settings, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching safety settings", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch safety settings")
	}

	if settings == nil {
		defaultSettings, err := model.NewSafetySettings(userID)
		if err != nil {
			slog.Error("Error creating default settings", "user_id", userID, "error", err)
			return nil, errors.New("failed to create default settings")
		}

		if err := uc.repo.Upsert(ctx, defaultSettings); err != nil {
			slog.Error("Error saving default settings", "user_id", userID, "error", err)
			return nil, errors.New("failed to save default settings")
		}

		settings = defaultSettings
	}

	return toResponse(settings), nil
}

func (uc *SafetySettingsUseCase) UpdateSettings(ctx context.Context, userID uuid.UUID, input dto.UpdateSafetySettingsInput) (*dto.SafetySettingsResponse, error) {
	settings, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching settings for update", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch settings")
	}

	if settings == nil {
		settings, err = model.NewSafetySettings(userID)
		if err != nil {
			slog.Error("Error creating settings", "user_id", userID, "error", err)
			return nil, errors.New("failed to create settings")
		}
	}

	if err := applyUpdates(settings, input); err != nil {
		return nil, err
	}

	if err := settings.Validate(); err != nil {
		return nil, err
	}

	settings.UpdatedAt = time.Now()

	if err := uc.repo.Upsert(ctx, settings); err != nil {
		slog.Error("Error updating settings", "user_id", userID, "error", err)
		return nil, errors.New("failed to update settings")
	}

	return toResponse(settings), nil
}

func applyUpdates(settings *model.SafetySettings, input dto.UpdateSafetySettingsInput) error {
	if input.NotificationsEnabled != nil {
		settings.NotificationsEnabled = *input.NotificationsEnabled
	}
	if input.NotificationAlertTypes != nil {
		settings.NotificationAlertTypes = *input.NotificationAlertTypes
	}
	if input.NotificationAlertRadiusMins != nil {
		settings.NotificationAlertRadiusMins = *input.NotificationAlertRadiusMins
	}
	if input.NotificationReportTypes != nil {
		settings.NotificationReportTypes = *input.NotificationReportTypes
	}
	if input.NotificationReportRadiusMins != nil {
		settings.NotificationReportRadiusMins = *input.NotificationReportRadiusMins
	}

	if input.LocationSharingEnabled != nil {
		settings.LocationSharingEnabled = *input.LocationSharingEnabled
	}
	if input.LocationHistoryEnabled != nil {
		settings.LocationHistoryEnabled = *input.LocationHistoryEnabled
	}

	if input.ProfileVisibility != nil {
		settings.ProfileVisibility = model.ProfileVisibility(*input.ProfileVisibility)
	}
	if input.AnonymousReports != nil {
		settings.AnonymousReports = *input.AnonymousReports
	}
	if input.ShowOnlineStatus != nil {
		settings.ShowOnlineStatus = *input.ShowOnlineStatus
	}

	if input.AutoAlertsEnabled != nil {
		settings.AutoAlertsEnabled = *input.AutoAlertsEnabled
	}
	if input.DangerZonesEnabled != nil {
		settings.DangerZonesEnabled = *input.DangerZonesEnabled
	}
	if input.TimeBasedAlertsEnabled != nil {
		settings.TimeBasedAlertsEnabled = *input.TimeBasedAlertsEnabled
	}

	if input.HighRiskStartTime != nil {
		t, err := time.Parse("15:04", *input.HighRiskStartTime)
		if err != nil {
			return errors.New("invalid high_risk_start_time format, expected HH:MM")
		}
		settings.HighRiskStartTime = t
	}
	if input.HighRiskEndTime != nil {
		t, err := time.Parse("15:04", *input.HighRiskEndTime)
		if err != nil {
			return errors.New("invalid high_risk_end_time format, expected HH:MM")
		}
		settings.HighRiskEndTime = t
	}

	if input.NightModeEnabled != nil {
		settings.NightModeEnabled = *input.NightModeEnabled
	}
	if input.NightModeStartTime != nil {
		t, err := time.Parse("15:04", *input.NightModeStartTime)
		if err != nil {
			return errors.New("invalid night_mode_start_time format, expected HH:MM")
		}
		settings.NightModeStartTime = t
	}
	if input.NightModeEndTime != nil {
		t, err := time.Parse("15:04", *input.NightModeEndTime)
		if err != nil {
			return errors.New("invalid night_mode_end_time format, expected HH:MM")
		}
		settings.NightModeEndTime = t
	}

	return nil
}

func toResponse(settings *model.SafetySettings) *dto.SafetySettingsResponse {
	return &dto.SafetySettingsResponse{
		ID:                           settings.ID.String(),
		UserID:                       settings.UserID.String(),
		NotificationsEnabled:         settings.NotificationsEnabled,
		NotificationAlertTypes:       settings.NotificationAlertTypes,
		NotificationAlertRadiusMins:  settings.NotificationAlertRadiusMins,
		NotificationReportTypes:      settings.NotificationReportTypes,
		NotificationReportRadiusMins: settings.NotificationReportRadiusMins,
		LocationSharingEnabled:       settings.LocationSharingEnabled,
		LocationHistoryEnabled:       settings.LocationHistoryEnabled,
		ProfileVisibility:            string(settings.ProfileVisibility),
		AnonymousReports:             settings.AnonymousReports,
		ShowOnlineStatus:             settings.ShowOnlineStatus,
		AutoAlertsEnabled:            settings.AutoAlertsEnabled,
		DangerZonesEnabled:           settings.DangerZonesEnabled,
		TimeBasedAlertsEnabled:       settings.TimeBasedAlertsEnabled,
		HighRiskStartTime:            settings.HighRiskStartTime.Format("15:04"),
		HighRiskEndTime:              settings.HighRiskEndTime.Format("15:04"),
		NightModeEnabled:             settings.NightModeEnabled,
		NightModeStartTime:           settings.NightModeStartTime.Format("15:04"),
		NightModeEndTime:             settings.NightModeEndTime.Format("15:04"),
		CreatedAt:                    settings.CreatedAt,
		UpdatedAt:                    settings.UpdatedAt,
	}
}
