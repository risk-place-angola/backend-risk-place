package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type SettingsChecker interface {
	CanReceiveNotifications(ctx context.Context, userID uuid.UUID, deviceID string) bool
	CanReceiveAlerts(ctx context.Context, userID uuid.UUID, deviceID string, severity string, distance int) bool
	CanReceiveReports(ctx context.Context, userID uuid.UUID, deviceID string, isVerified bool, distance int) bool
	CanShareLocation(ctx context.Context, userID uuid.UUID, deviceID string) bool
	CanSaveLocationHistory(ctx context.Context, userID uuid.UUID, deviceID string) bool
	ShouldShowOnline(ctx context.Context, userID uuid.UUID) bool
	IsInHighRiskTime(ctx context.Context, userID uuid.UUID, deviceID string) bool
	HasDangerZonesEnabled(ctx context.Context, userID uuid.UUID, deviceID string) bool
}

type settingsChecker struct {
	settingsRepo         repository.SafetySettingsRepository
	anonymousSessionRepo repository.AnonymousSessionRepository
}

func NewSettingsChecker(
	settingsRepo repository.SafetySettingsRepository,
	anonymousSessionRepo repository.AnonymousSessionRepository,
) SettingsChecker {
	return &settingsChecker{
		settingsRepo:         settingsRepo,
		anonymousSessionRepo: anonymousSessionRepo,
	}
}

func (s *settingsChecker) CanReceiveNotifications(ctx context.Context, userID uuid.UUID, deviceID string) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return true
	}
	return settings.NotificationsEnabled
}

func (s *settingsChecker) CanReceiveAlerts(ctx context.Context, userID uuid.UUID, deviceID string, severity string, distance int) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return true
	}

	if !settings.NotificationsEnabled {
		return false
	}

	if distance > settings.NotificationAlertRadiusMins {
		return false
	}

	for _, allowedSeverity := range settings.NotificationAlertTypes {
		if allowedSeverity == severity {
			return true
		}
	}

	return false
}

func (s *settingsChecker) CanReceiveReports(ctx context.Context, userID uuid.UUID, deviceID string, isVerified bool, distance int) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return true
	}

	if !settings.NotificationsEnabled {
		return false
	}

	if distance > settings.NotificationReportRadiusMins {
		return false
	}

	for _, reportType := range settings.NotificationReportTypes {
		if reportType == "all" || (reportType == "verified" && isVerified) {
			return true
		}
	}

	return false
}

func (s *settingsChecker) CanShareLocation(ctx context.Context, userID uuid.UUID, deviceID string) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return false
	}
	return settings.LocationSharingEnabled
}

func (s *settingsChecker) CanSaveLocationHistory(ctx context.Context, userID uuid.UUID, deviceID string) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return true
	}
	return settings.LocationHistoryEnabled
}

func (s *settingsChecker) ShouldShowOnline(ctx context.Context, userID uuid.UUID) bool {
	if userID == uuid.Nil {
		return false
	}

	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil || settings == nil {
		return true
	}
	return settings.ShowOnlineStatus
}

func (s *settingsChecker) IsInHighRiskTime(ctx context.Context, userID uuid.UUID, deviceID string) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return false
	}

	if !settings.TimeBasedAlertsEnabled {
		return false
	}

	now := time.Now().UTC()
	currentTime := now.Format("15:04")

	startTime := settings.HighRiskStartTime.UTC().Format("15:04")
	endTime := settings.HighRiskEndTime.UTC().Format("15:04")

	if startTime > endTime {
		return currentTime >= startTime || currentTime <= endTime
	}

	return currentTime >= startTime && currentTime <= endTime
}

func (s *settingsChecker) HasDangerZonesEnabled(ctx context.Context, userID uuid.UUID, deviceID string) bool {
	settings, err := s.getSettings(ctx, userID, deviceID)
	if err != nil || settings == nil {
		return true
	}
	return settings.DangerZonesEnabled
}

var ErrNoIdentifier = errors.New("no user ID or device ID provided")

func (s *settingsChecker) getSettings(ctx context.Context, userID uuid.UUID, deviceID string) (*model.SafetySettings, error) {
	if userID != uuid.Nil {
		settings, err := s.settingsRepo.GetByUserID(ctx, userID)
		if err != nil {
			slog.Debug("failed to get settings by user ID", "user_id", userID, "error", err)
			return nil, err
		}
		return settings, nil
	}

	if deviceID != "" {
		settings, err := s.settingsRepo.GetByDeviceID(ctx, deviceID)
		if err != nil {
			slog.Debug("failed to get settings by device ID", "device_id", deviceID, "error", err)
			return nil, err
		}
		return settings, nil
	}

	return nil, ErrNoIdentifier
}
