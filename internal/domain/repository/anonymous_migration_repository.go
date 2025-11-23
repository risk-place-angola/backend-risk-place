package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type DeviceUserMappingRepository interface {
	Create(ctx context.Context, mapping *model.DeviceUserMapping) error
	GetActiveMapping(ctx context.Context, deviceID string) (*model.DeviceUserMapping, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.DeviceUserMapping, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type AnonymousMigrationRepository interface {
	Create(ctx context.Context, migration *model.AnonymousUserMigration) error
	UpdateCounters(ctx context.Context, id uuid.UUID, alertsMigrated, subscriptionsMigrated, locationSharingsMigrated int, settingsMigrated bool) error
	MarkCompleted(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errorMessage string) error
	GetByDeviceID(ctx context.Context, deviceID string) ([]*model.AnonymousUserMigration, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.AnonymousUserMigration, error)
}
