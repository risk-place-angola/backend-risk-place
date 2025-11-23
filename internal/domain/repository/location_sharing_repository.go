package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type LocationSharingRepository interface {
	GenericRepository[model.LocationSharing]
	FindByToken(ctx context.Context, token string) (*model.LocationSharing, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*model.LocationSharing, error)
	FindActiveByDeviceID(ctx context.Context, deviceID string) ([]*model.LocationSharing, error)
	DeactivateExpired(ctx context.Context) error
}
