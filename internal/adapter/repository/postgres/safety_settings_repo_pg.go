package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type safetySettingsRepoPG struct {
	q sqlc.Querier
}

func NewSafetySettingsRepository(db *sql.DB) repository.SafetySettingsRepository {
	return &safetySettingsRepoPG{q: sqlc.New(db)}
}

func (r *safetySettingsRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.SafetySettings, error) {
	row, err := r.q.GetSafetySettingsByUserID(ctx, uuidToNullUUID(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *safetySettingsRepoPG) Upsert(ctx context.Context, settings *model.SafetySettings) error {
	return r.q.UpsertSafetySettings(ctx, sqlc.UpsertSafetySettingsParams{
		ID:     settings.ID,
		UserID: uuidPtrToNullUUID(settings.UserID),
		NotificationsEnabled: sql.NullBool{
			Bool:  settings.NotificationsEnabled,
			Valid: true,
		},
		NotificationAlertTypes: settings.NotificationAlertTypes,
		NotificationAlertRadiusMins: sql.NullInt32{
			Int32: safeIntToInt32(settings.NotificationAlertRadiusMins),
			Valid: true,
		},
		NotificationReportTypes: settings.NotificationReportTypes,
		NotificationReportRadiusMins: sql.NullInt32{
			Int32: safeIntToInt32(settings.NotificationReportRadiusMins),
			Valid: true,
		},
		LocationSharingEnabled: sql.NullBool{
			Bool:  settings.LocationSharingEnabled,
			Valid: true,
		},
		LocationHistoryEnabled: sql.NullBool{
			Bool:  settings.LocationHistoryEnabled,
			Valid: true,
		},
		ProfileVisibility: sql.NullString{
			String: string(settings.ProfileVisibility),
			Valid:  true,
		},
		AnonymousReports: sql.NullBool{
			Bool:  settings.AnonymousReports,
			Valid: true,
		},
		ShowOnlineStatus: sql.NullBool{
			Bool:  settings.ShowOnlineStatus,
			Valid: true,
		},
		AutoAlertsEnabled: sql.NullBool{
			Bool:  settings.AutoAlertsEnabled,
			Valid: true,
		},
		DangerZonesEnabled: sql.NullBool{
			Bool:  settings.DangerZonesEnabled,
			Valid: true,
		},
		TimeBasedAlertsEnabled: sql.NullBool{
			Bool:  settings.TimeBasedAlertsEnabled,
			Valid: true,
		},
		HighRiskStartTime: sql.NullTime{
			Time:  settings.HighRiskStartTime,
			Valid: true,
		},
		HighRiskEndTime: sql.NullTime{
			Time:  settings.HighRiskEndTime,
			Valid: true,
		},
		NightModeEnabled: sql.NullBool{
			Bool:  settings.NightModeEnabled,
			Valid: true,
		},
		NightModeStartTime: sql.NullTime{
			Time:  settings.NightModeStartTime,
			Valid: true,
		},
		NightModeEndTime: sql.NullTime{
			Time:  settings.NightModeEndTime,
			Valid: true,
		},
		CreatedAt: sql.NullTime{
			Time:  settings.CreatedAt,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (r *safetySettingsRepoPG) GetByDeviceID(ctx context.Context, deviceID string) (*model.SafetySettings, error) {
	row, err := r.q.GetSafetySettingsByAnonymousSessionID(ctx, sqlc.GetSafetySettingsByAnonymousSessionIDParams{
		DeviceID: sql.NullString{String: deviceID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *safetySettingsRepoPG) UpsertAnonymous(ctx context.Context, settings *model.SafetySettings) error {
	// First, try to get existing record by device_id to reuse its ID
	existing, err := r.GetByDeviceID(ctx, *settings.DeviceID)
	settingsID := settings.ID
	if err == nil && existing != nil {
		settingsID = existing.ID
	}

	return r.q.UpsertAnonymousSafetySettings(ctx, sqlc.UpsertAnonymousSafetySettingsParams{
		ID:                 settingsID,
		AnonymousSessionID: uuidPtrToNullUUID(settings.AnonymousSessionID),
		DeviceID:           sql.NullString{String: *settings.DeviceID, Valid: true},
		NotificationsEnabled: sql.NullBool{
			Bool:  settings.NotificationsEnabled,
			Valid: true,
		},
		NotificationAlertTypes: settings.NotificationAlertTypes,
		NotificationAlertRadiusMins: sql.NullInt32{
			Int32: safeIntToInt32(settings.NotificationAlertRadiusMins),
			Valid: true,
		},
		NotificationReportTypes: settings.NotificationReportTypes,
		NotificationReportRadiusMins: sql.NullInt32{
			Int32: safeIntToInt32(settings.NotificationReportRadiusMins),
			Valid: true,
		},
		LocationSharingEnabled: sql.NullBool{
			Bool:  settings.LocationSharingEnabled,
			Valid: true,
		},
		LocationHistoryEnabled: sql.NullBool{
			Bool:  settings.LocationHistoryEnabled,
			Valid: true,
		},
		ProfileVisibility: sql.NullString{
			String: string(settings.ProfileVisibility),
			Valid:  true,
		},
		AnonymousReports: sql.NullBool{
			Bool:  settings.AnonymousReports,
			Valid: true,
		},
		ShowOnlineStatus: sql.NullBool{
			Bool:  settings.ShowOnlineStatus,
			Valid: true,
		},
		AutoAlertsEnabled: sql.NullBool{
			Bool:  settings.AutoAlertsEnabled,
			Valid: true,
		},
		DangerZonesEnabled: sql.NullBool{
			Bool:  settings.DangerZonesEnabled,
			Valid: true,
		},
		TimeBasedAlertsEnabled: sql.NullBool{
			Bool:  settings.TimeBasedAlertsEnabled,
			Valid: true,
		},
		HighRiskStartTime: sql.NullTime{
			Time:  settings.HighRiskStartTime,
			Valid: true,
		},
		HighRiskEndTime: sql.NullTime{
			Time:  settings.HighRiskEndTime,
			Valid: true,
		},
		NightModeEnabled: sql.NullBool{
			Bool:  settings.NightModeEnabled,
			Valid: true,
		},
		NightModeStartTime: sql.NullTime{
			Time:  settings.NightModeStartTime,
			Valid: true,
		},
		NightModeEndTime: sql.NullTime{
			Time:  settings.NightModeEndTime,
			Valid: true,
		},
		CreatedAt: sql.NullTime{
			Time:  settings.CreatedAt,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (r *safetySettingsRepoPG) toDomain(row sqlc.UserSafetySetting) *model.SafetySettings {
	return &model.SafetySettings{
		ID:                           row.ID,
		UserID:                       nullUUIDToPtr(row.UserID),
		AnonymousSessionID:           nullUUIDToPtr(row.AnonymousSessionID),
		DeviceID:                     nullStringToPtr(row.DeviceID),
		NotificationsEnabled:         row.NotificationsEnabled.Bool,
		NotificationAlertTypes:       row.NotificationAlertTypes,
		NotificationAlertRadiusMins:  int(row.NotificationAlertRadiusMins.Int32),
		NotificationReportTypes:      row.NotificationReportTypes,
		NotificationReportRadiusMins: int(row.NotificationReportRadiusMins.Int32),
		LocationSharingEnabled:       row.LocationSharingEnabled.Bool,
		LocationHistoryEnabled:       row.LocationHistoryEnabled.Bool,
		ProfileVisibility:            model.ProfileVisibility(row.ProfileVisibility.String),
		AnonymousReports:             row.AnonymousReports.Bool,
		ShowOnlineStatus:             row.ShowOnlineStatus.Bool,
		AutoAlertsEnabled:            row.AutoAlertsEnabled.Bool,
		DangerZonesEnabled:           row.DangerZonesEnabled.Bool,
		TimeBasedAlertsEnabled:       row.TimeBasedAlertsEnabled.Bool,
		HighRiskStartTime:            row.HighRiskStartTime.Time,
		HighRiskEndTime:              row.HighRiskEndTime.Time,
		NightModeEnabled:             row.NightModeEnabled.Bool,
		NightModeStartTime:           row.NightModeStartTime.Time,
		NightModeEndTime:             row.NightModeEndTime.Time,
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
}
