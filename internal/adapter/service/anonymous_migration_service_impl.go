package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type anonymousMigrationService struct {
	deviceMappingRepo    repository.DeviceUserMappingRepository
	migrationRepo        repository.AnonymousMigrationRepository
	anonymousSessionRepo repository.AnonymousSessionRepository
	alertRepo            repository.AlertRepository
	safetySettingsRepo   repository.SafetySettingsRepository
	locationSharingRepo  repository.LocationSharingRepository
}

func NewAnonymousMigrationService(
	deviceMappingRepo repository.DeviceUserMappingRepository,
	migrationRepo repository.AnonymousMigrationRepository,
	anonymousSessionRepo repository.AnonymousSessionRepository,
	alertRepo repository.AlertRepository,
	safetySettingsRepo repository.SafetySettingsRepository,
	locationSharingRepo repository.LocationSharingRepository,
) domainService.AnonymousMigrationService {
	return &anonymousMigrationService{
		deviceMappingRepo:    deviceMappingRepo,
		migrationRepo:        migrationRepo,
		anonymousSessionRepo: anonymousSessionRepo,
		alertRepo:            alertRepo,
		safetySettingsRepo:   safetySettingsRepo,
		locationSharingRepo:  locationSharingRepo,
	}
}

func (s *anonymousMigrationService) MigrateAnonymousData(
	ctx context.Context,
	deviceID string,
	userID uuid.UUID,
	migrationType string,
) error {
	slog.Info("Starting anonymous data migration",
		"device_id", deviceID,
		"user_id", userID,
		"migration_type", migrationType)

	existingMapping, err := s.deviceMappingRepo.GetActiveMapping(ctx, deviceID)
	if err != nil && !errors.Is(err, domainErrors.ErrNotFound) {
		slog.Error("Failed to check existing mapping", "device_id", deviceID, "error", err)
		return fmt.Errorf("failed to check existing mapping: %w", err)
	}

	if existingMapping != nil {
		if existingMapping.UserID == userID {
			slog.Info("Device already mapped to this user", "device_id", deviceID, "user_id", userID)
			return nil
		}

		if err := s.deviceMappingRepo.Deactivate(ctx, existingMapping.ID); err != nil {
			slog.Error("Failed to deactivate old mapping", "device_id", deviceID, "error", err)
			return fmt.Errorf("failed to deactivate old mapping: %w", err)
		}
	}

	anonymousSession, err := s.anonymousSessionRepo.FindByDeviceID(ctx, deviceID)
	if err != nil && !errors.Is(err, domainErrors.ErrNotFound) {
		slog.Error("Failed to fetch anonymous session", "device_id", deviceID, "error", err)
		return fmt.Errorf("failed to fetch anonymous session: %w", err)
	}

	if anonymousSession == nil {
		slog.Debug("No anonymous session found, creating device mapping only", "device_id", deviceID)

		mapping, err := model.NewDeviceUserMapping(deviceID, uuid.Nil, userID)
		if err != nil {
			return fmt.Errorf("failed to create device mapping: %w", err)
		}

		if err := s.deviceMappingRepo.Create(ctx, mapping); err != nil {
			return fmt.Errorf("failed to save device mapping: %w", err)
		}

		slog.Info("Device mapped to user successfully", "device_id", deviceID, "user_id", userID)
		return nil
	}

	// 2. Criar log de migração
	migration, err := model.NewAnonymousUserMigration(
		anonymousSession.ID,
		deviceID,
		userID,
		migrationType,
	)
	if err != nil {
		slog.Error("Failed to create migration log", "error", err)
		return fmt.Errorf("failed to create migration log: %w", err)
	}

	if err := s.migrationRepo.Create(ctx, migration); err != nil {
		slog.Error("Failed to save migration log", "error", err)
		return fmt.Errorf("failed to save migration log: %w", err)
	}

	alertsMigrated := 0
	subscriptionsMigrated := 0
	locationSharingsMigrated := 0
	settingsMigrated := false

	if err := s.migrationRepo.UpdateCounters(
		ctx,
		migration.ID,
		alertsMigrated,
		subscriptionsMigrated,
		locationSharingsMigrated,
		settingsMigrated,
	); err != nil {
		slog.Error("Failed to update migration counters", "error", err)
		return fmt.Errorf("failed to update migration counters: %w", err)
	}

	// 5. Criar device mapping
	mapping, err := model.NewDeviceUserMapping(deviceID, anonymousSession.ID, userID)
	if err != nil {
		slog.Error("Failed to create device mapping", "error", err)
		if markErr := s.migrationRepo.MarkFailed(ctx, migration.ID, err.Error()); markErr != nil {
			slog.Error("Failed to mark migration as failed", "error", markErr)
		}
		return fmt.Errorf("failed to create device mapping: %w", err)
	}

	if err := s.deviceMappingRepo.Create(ctx, mapping); err != nil {
		slog.Error("Failed to save device mapping", "error", err)
		if markErr := s.migrationRepo.MarkFailed(ctx, migration.ID, err.Error()); markErr != nil {
			slog.Error("Failed to mark migration as failed", "error", markErr)
		}
		return fmt.Errorf("failed to save device mapping: %w", err)
	}

	if err := s.migrationRepo.MarkCompleted(ctx, migration.ID); err != nil {
		slog.Error("Failed to mark migration as completed", "error", err)
		return fmt.Errorf("failed to mark migration as completed: %w", err)
	}

	slog.Info("Anonymous data migration completed successfully",
		"migration_id", migration.ID,
		"alerts_migrated", alertsMigrated,
		"subscriptions_migrated", subscriptionsMigrated,
		"settings_migrated", settingsMigrated)

	return nil
}

func (s *anonymousMigrationService) CheckExistingMapping(ctx context.Context, deviceID string) (*uuid.UUID, error) {
	mapping, err := s.deviceMappingRepo.GetActiveMapping(ctx, deviceID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		slog.Error("Failed to check existing mapping", "device_id", deviceID, "error", err)
		return nil, fmt.Errorf("failed to check existing mapping: %w", err)
	}

	return &mapping.UserID, nil
}

func (s *anonymousMigrationService) RollbackMigration(ctx context.Context, migrationID uuid.UUID) error {
	slog.Warn("Rollback not implemented yet", "migration_id", migrationID)
	return fmt.Errorf("rollback not implemented yet")
}
