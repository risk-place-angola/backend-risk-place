package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type deviceUserMappingRepoPG struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewDeviceUserMappingRepository creates a new DeviceUserMappingRepository implementation
func NewDeviceUserMappingRepository(db *sql.DB) repository.DeviceUserMappingRepository {
	return &deviceUserMappingRepoPG{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *deviceUserMappingRepoPG) Create(ctx context.Context, mapping *model.DeviceUserMapping) error {
	if err := mapping.Validate(); err != nil {
		return fmt.Errorf("invalid device user mapping: %w", err)
	}

	err := r.queries.CreateDeviceUserMapping(ctx, sqlc.CreateDeviceUserMappingParams{
		ID:                 mapping.ID,
		DeviceID:           mapping.DeviceID,
		AnonymousSessionID: mapping.AnonymousSessionID,
		UserID:             mapping.UserID,
	})

	if err != nil {
		return fmt.Errorf("failed to create device user mapping: %w", err)
	}

	return nil
}

func (r *deviceUserMappingRepoPG) GetActiveMapping(ctx context.Context, deviceID string) (*model.DeviceUserMapping, error) {
	sqlcMapping, err := r.queries.GetActiveDeviceUserMapping(ctx, deviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get active device user mapping: %w", err)
	}

	return r.sqlcToModel(&sqlcMapping), nil
}

func (r *deviceUserMappingRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.DeviceUserMapping, error) {
	sqlcMappings, err := r.queries.GetDeviceUserMappingsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device user mappings by user id: %w", err)
	}

	mappings := make([]*model.DeviceUserMapping, 0, len(sqlcMappings))
	for _, sqlcMapping := range sqlcMappings {
		mappings = append(mappings, r.sqlcToModel(&sqlcMapping))
	}

	return mappings, nil
}

func (r *deviceUserMappingRepoPG) Deactivate(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeactivateDeviceUserMapping(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate device user mapping: %w", err)
	}

	return nil
}

// sqlcToModel converts SQLC model to domain model
func (r *deviceUserMappingRepoPG) sqlcToModel(sqlcMapping *sqlc.DeviceUserMapping) *model.DeviceUserMapping {
	mapping := &model.DeviceUserMapping{
		ID:                 sqlcMapping.ID,
		DeviceID:           sqlcMapping.DeviceID,
		AnonymousSessionID: sqlcMapping.AnonymousSessionID,
		UserID:             sqlcMapping.UserID,
		MappedAt:           sqlcMapping.MappedAt,
		IsActive:           sqlcMapping.IsActive,
	}

	if sqlcMapping.UnmappedAt.Valid {
		mapping.UnmappedAt = &sqlcMapping.UnmappedAt.Time
	}

	return mapping
}
