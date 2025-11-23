package service

import (
	"context"

	"github.com/google/uuid"
)

type AnonymousMigrationService interface {
	MigrateAnonymousData(ctx context.Context, deviceID string, userID uuid.UUID, migrationType string) error
	CheckExistingMapping(ctx context.Context, deviceID string) (*uuid.UUID, error)
	RollbackMigration(ctx context.Context, migrationID uuid.UUID) error
}
