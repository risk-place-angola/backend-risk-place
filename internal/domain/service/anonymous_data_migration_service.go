package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type AnonymousDataMigrationService struct {
	anonymousSessionRepo AnonymousSessionRepository
	alertRepo            AlertRepository
	alertSubRepo         AlertSubscriptionRepository
	safetySettingsRepo   SafetySettingsRepository
	locationSharingRepo  LocationSharingRepository
	migrationLogRepo     MigrationLogRepository
}

func NewAnonymousDataMigrationService(
	anonymousSessionRepo AnonymousSessionRepository,
	alertRepo AlertRepository,
	alertSubRepo AlertSubscriptionRepository,
	safetySettingsRepo SafetySettingsRepository,
	locationSharingRepo LocationSharingRepository,
	migrationLogRepo MigrationLogRepository,
) *AnonymousDataMigrationService {
	return &AnonymousDataMigrationService{
		anonymousSessionRepo: anonymousSessionRepo,
		alertRepo:            alertRepo,
		alertSubRepo:         alertSubRepo,
		safetySettingsRepo:   safetySettingsRepo,
		locationSharingRepo:  locationSharingRepo,
		migrationLogRepo:     migrationLogRepo,
	}
}

type MigrationResult struct {
	AlertsMigrated              int
	SubscriptionsMigrated       int
	SettingsMigrated            bool
	LocationSharingsMigrated    int
	AnonymousSessionDeactivated bool
	Error                       error
}

type MigrationType string

const (
	MigrationTypeSignup MigrationType = "signup"
	MigrationTypeLogin  MigrationType = "login"
	MigrationTypeManual MigrationType = "manual"
)

func (s *AnonymousDataMigrationService) MigrateAnonymousUser(
	ctx context.Context,
	deviceID string,
	userID uuid.UUID,
	migrationType MigrationType,
) (*MigrationResult, error) {
	result := &MigrationResult{}

	session, err := s.anonymousSessionRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, errors.New("anonymous session not found")
	}

	if session.MigratedToUserID != nil {
		return nil, errors.New("device already migrated to another user")
	}

	migrationLog, err := s.migrationLogRepo.Create(ctx, &model.AnonymousUserMigration{
		ID:                 uuid.New(),
		AnonymousSessionID: session.ID,
		DeviceID:           deviceID,
		UserID:             userID,
		MigrationType:      string(migrationType),
		StartedAt:          time.Now(),
	})
	if err != nil {
		slog.Error("Failed to create migration log", "error", err)
		return nil, err
	}

	tx, err := s.beginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer s.rollbackIfNeeded(tx)

	if err := s.migrateAlerts(ctx, session.ID, userID, result); err != nil {
		migrationLog.FailedAt = timePtr(time.Now())
		migrationLog.ErrorMessage = strPtr(err.Error())
		_ = s.migrationLogRepo.Update(ctx, migrationLog)
		return nil, err
	}

	if err := s.migrateSubscriptions(ctx, session.ID, userID, result); err != nil {
		migrationLog.FailedAt = timePtr(time.Now())
		migrationLog.ErrorMessage = strPtr(err.Error())
		_ = s.migrationLogRepo.Update(ctx, migrationLog)
		return nil, err
	}

	if err := s.migrateSettings(ctx, session.ID, userID, result); err != nil {
		migrationLog.FailedAt = timePtr(time.Now())
		migrationLog.ErrorMessage = strPtr(err.Error())
		_ = s.migrationLogRepo.Update(ctx, migrationLog)
		return nil, err
	}

	if err := s.migrateLocationSharings(ctx, session.ID, userID, result); err != nil {
		migrationLog.FailedAt = timePtr(time.Now())
		migrationLog.ErrorMessage = strPtr(err.Error())
		_ = s.migrationLogRepo.Update(ctx, migrationLog)
		return nil, err
	}

	session.MigratedToUserID = &userID
	session.MigratedAt = timePtr(time.Now())
	session.IsActive = false

	if err := s.anonymousSessionRepo.Update(ctx, session); err != nil {
		migrationLog.FailedAt = timePtr(time.Now())
		migrationLog.ErrorMessage = strPtr("failed to deactivate anonymous session")
		_ = s.migrationLogRepo.Update(ctx, migrationLog)
		return nil, err
	}

	result.AnonymousSessionDeactivated = true

	if err := s.createDeviceUserMapping(ctx, deviceID, session.ID, userID); err != nil {
		slog.Warn("Failed to create device-user mapping", "error", err)
	}

	if err := s.commitTransaction(tx); err != nil {
		migrationLog.FailedAt = timePtr(time.Now())
		migrationLog.ErrorMessage = strPtr("failed to commit transaction")
		_ = s.migrationLogRepo.Update(ctx, migrationLog)
		return nil, err
	}

	migrationLog.AlertsMigrated = result.AlertsMigrated
	migrationLog.SubscriptionsMigrated = result.SubscriptionsMigrated
	migrationLog.SettingsMigrated = result.SettingsMigrated
	migrationLog.LocationSharingsMigrated = result.LocationSharingsMigrated
	migrationLog.CompletedAt = timePtr(time.Now())

	if err := s.migrationLogRepo.Update(ctx, migrationLog); err != nil {
		slog.Error("Failed to update migration log", "error", err)
	}

	slog.Info("Anonymous user migration completed",
		"device_id", deviceID,
		"user_id", userID,
		"alerts", result.AlertsMigrated,
		"subscriptions", result.SubscriptionsMigrated,
		"settings", result.SettingsMigrated,
	)

	return result, nil
}

func (s *AnonymousDataMigrationService) migrateAlerts(
	ctx context.Context,
	anonymousSessionID, userID uuid.UUID,
	result *MigrationResult,
) error {
	alerts, err := s.alertRepo.GetByAnonymousSessionID(ctx, anonymousSessionID)
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		alert.CreatedBy = &userID
		alert.AnonymousSessionID = nil
		alert.DeviceID = nil

		if err := s.alertRepo.Update(ctx, alert); err != nil {
			return err
		}
		result.AlertsMigrated++
	}

	return nil
}

func (s *AnonymousDataMigrationService) migrateSubscriptions(
	ctx context.Context,
	anonymousSessionID, userID uuid.UUID,
	result *MigrationResult,
) error {
	subscriptions, err := s.alertSubRepo.GetByAnonymousSessionID(ctx, anonymousSessionID)
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		sub.UserID = &userID
		sub.AnonymousSessionID = nil
		sub.DeviceID = nil

		if err := s.alertSubRepo.Update(ctx, sub); err != nil {
			return err
		}
		result.SubscriptionsMigrated++
	}

	return nil
}

func (s *AnonymousDataMigrationService) migrateSettings(
	ctx context.Context,
	anonymousSessionID, userID uuid.UUID,
	result *MigrationResult,
) error {
	settings, err := s.safetySettingsRepo.GetByAnonymousSessionID(ctx, anonymousSessionID)
	if err != nil {
		return err
	}
	if settings == nil {
		return nil
	}

	existingSettings, err := s.safetySettingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existingSettings == nil {
		settings.UserID = &userID
		settings.AnonymousSessionID = nil
		settings.DeviceID = nil

		if err := s.safetySettingsRepo.Upsert(ctx, settings); err != nil {
			return err
		}
		result.SettingsMigrated = true
	} else {
		if err := s.safetySettingsRepo.DeleteByAnonymousSessionID(ctx, anonymousSessionID); err != nil {
			slog.Warn("Failed to delete anonymous settings", "error", err)
		}
	}

	return nil
}

func (s *AnonymousDataMigrationService) migrateLocationSharings(
	ctx context.Context,
	anonymousSessionID, userID uuid.UUID,
	result *MigrationResult,
) error {
	sharings, err := s.locationSharingRepo.GetByAnonymousSessionID(ctx, anonymousSessionID)
	if err != nil {
		return err
	}

	for _, sharing := range sharings {
		sharing.UserID = &userID
		sharing.AnonymousSessionID = nil
		sharing.DeviceID = nil

		if err := s.locationSharingRepo.Update(ctx, sharing); err != nil {
			return err
		}
		result.LocationSharingsMigrated++
	}

	return nil
}

func (s *AnonymousDataMigrationService) createDeviceUserMapping(
	ctx context.Context,
	deviceID string,
	anonymousSessionID, userID uuid.UUID,
) error {
	mapping := &model.DeviceUserMapping{
		ID:                 uuid.New(),
		DeviceID:           deviceID,
		AnonymousSessionID: anonymousSessionID,
		UserID:             userID,
		MappedAt:           time.Now(),
		IsActive:           true,
	}

	return s.migrationLogRepo.CreateDeviceUserMapping(ctx, mapping)
}

func (s *AnonymousDataMigrationService) beginTransaction(ctx context.Context) (Transaction, error) {
	return nil, fmt.Errorf("transaction not implemented")
}

func (s *AnonymousDataMigrationService) commitTransaction(tx Transaction) error {
	return nil
}

func (s *AnonymousDataMigrationService) rollbackIfNeeded(tx Transaction) {
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func strPtr(s string) *string {
	return &s
}

type Transaction interface{}

type AnonymousSessionRepository interface {
	GetByDeviceID(ctx context.Context, deviceID string) (*model.AnonymousSession, error)
	Update(ctx context.Context, session *model.AnonymousSession) error
}

type AlertRepository interface {
	GetByAnonymousSessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.Alert, error)
	Update(ctx context.Context, alert *model.Alert) error
}

type AlertSubscriptionRepository interface {
	GetByAnonymousSessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.AlertSubscription, error)
	Update(ctx context.Context, subscription *model.AlertSubscription) error
}

type SafetySettingsRepository interface {
	GetByAnonymousSessionID(ctx context.Context, sessionID uuid.UUID) (*model.SafetySettings, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.SafetySettings, error)
	Upsert(ctx context.Context, settings *model.SafetySettings) error
	DeleteByAnonymousSessionID(ctx context.Context, sessionID uuid.UUID) error
}

type LocationSharingRepository interface {
	GetByAnonymousSessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.LocationSharing, error)
	Update(ctx context.Context, sharing *model.LocationSharing) error
}

type MigrationLogRepository interface {
	Create(ctx context.Context, migration *model.AnonymousUserMigration) (*model.AnonymousUserMigration, error)
	Update(ctx context.Context, migration *model.AnonymousUserMigration) error
	CreateDeviceUserMapping(ctx context.Context, mapping *model.DeviceUserMapping) error
}
