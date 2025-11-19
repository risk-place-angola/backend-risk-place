package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

func safeIntToInt32(val int) int32 {
	if val > math.MaxInt32 {
		return math.MaxInt32
	}
	if val < math.MinInt32 {
		return math.MinInt32
	}
	return int32(val)
}

type anonymousMigrationRepoPG struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewAnonymousMigrationRepository creates a new AnonymousMigrationRepository implementation
func NewAnonymousMigrationRepository(db *sql.DB) repository.AnonymousMigrationRepository {
	return &anonymousMigrationRepoPG{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *anonymousMigrationRepoPG) Create(ctx context.Context, migration *model.AnonymousUserMigration) error {
	if err := migration.Validate(); err != nil {
		return fmt.Errorf("invalid anonymous user migration: %w", err)
	}

	err := r.queries.CreateAnonymousMigration(ctx, sqlc.CreateAnonymousMigrationParams{
		ID:                 migration.ID,
		AnonymousSessionID: migration.AnonymousSessionID,
		DeviceID:           migration.DeviceID,
		UserID:             migration.UserID,
		MigrationType:      migration.MigrationType,
	})

	if err != nil {
		return fmt.Errorf("failed to create anonymous user migration: %w", err)
	}

	return nil
}

func (r *anonymousMigrationRepoPG) UpdateCounters(
	ctx context.Context,
	id uuid.UUID,
	alertsMigrated, subscriptionsMigrated, locationSharingsMigrated int,
	settingsMigrated bool,
) error {
	err := r.queries.UpdateMigrationCounters(ctx, sqlc.UpdateMigrationCountersParams{
		ID:                       id,
		AlertsMigrated:           safeIntToInt32(alertsMigrated),
		SubscriptionsMigrated:    safeIntToInt32(subscriptionsMigrated),
		SettingsMigrated:         settingsMigrated,
		LocationSharingsMigrated: safeIntToInt32(locationSharingsMigrated),
	})

	if err != nil {
		return fmt.Errorf("failed to update migration counters: %w", err)
	}

	return nil
}

func (r *anonymousMigrationRepoPG) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	err := r.queries.MarkMigrationCompleted(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to mark migration as completed: %w", err)
	}

	return nil
}

func (r *anonymousMigrationRepoPG) MarkFailed(ctx context.Context, id uuid.UUID, errorMessage string) error {
	err := r.queries.MarkMigrationFailed(ctx, sqlc.MarkMigrationFailedParams{
		ID: id,
		ErrorMessage: sql.NullString{
			String: errorMessage,
			Valid:  true,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to mark migration as failed: %w", err)
	}

	return nil
}

func (r *anonymousMigrationRepoPG) GetByDeviceID(ctx context.Context, deviceID string) ([]*model.AnonymousUserMigration, error) {
	sqlcMigrations, err := r.queries.GetMigrationsByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get migrations by device id: %w", err)
	}

	migrations := make([]*model.AnonymousUserMigration, 0, len(sqlcMigrations))
	for _, sqlcMigration := range sqlcMigrations {
		migrations = append(migrations, r.sqlcToModel(&sqlcMigration))
	}

	return migrations, nil
}

func (r *anonymousMigrationRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.AnonymousUserMigration, error) {
	sqlcMigrations, err := r.queries.GetMigrationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get migrations by user id: %w", err)
	}

	migrations := make([]*model.AnonymousUserMigration, 0, len(sqlcMigrations))
	for _, sqlcMigration := range sqlcMigrations {
		migrations = append(migrations, r.sqlcToModel(&sqlcMigration))
	}

	return migrations, nil
}

// sqlcToModel converts SQLC model to domain model
func (r *anonymousMigrationRepoPG) sqlcToModel(sqlcMigration *sqlc.AnonymousUserMigration) *model.AnonymousUserMigration {
	migration := &model.AnonymousUserMigration{
		ID:                       sqlcMigration.ID,
		AnonymousSessionID:       sqlcMigration.AnonymousSessionID,
		DeviceID:                 sqlcMigration.DeviceID,
		UserID:                   sqlcMigration.UserID,
		AlertsMigrated:           int(sqlcMigration.AlertsMigrated),
		SubscriptionsMigrated:    int(sqlcMigration.SubscriptionsMigrated),
		SettingsMigrated:         sqlcMigration.SettingsMigrated,
		LocationSharingsMigrated: int(sqlcMigration.LocationSharingsMigrated),
		MigrationType:            sqlcMigration.MigrationType,
		StartedAt:                sqlcMigration.StartedAt,
	}

	if sqlcMigration.CompletedAt.Valid {
		migration.CompletedAt = &sqlcMigration.CompletedAt.Time
	}

	if sqlcMigration.FailedAt.Valid {
		migration.FailedAt = &sqlcMigration.FailedAt.Time
	}

	if sqlcMigration.ErrorMessage.Valid {
		errorMsg := sqlcMigration.ErrorMessage.String
		migration.ErrorMessage = &errorMsg
	}

	return migration
}
