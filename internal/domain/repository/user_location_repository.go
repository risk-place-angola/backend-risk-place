package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type UserLocationRepository interface {
	Upsert(ctx context.Context, location *model.UserLocation) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.UserLocation, error)
	FindNearbyUsers(ctx context.Context, lat, lon, radiusMeters float64, limit int) ([]*model.UserLocation, error)
	DeleteStale(ctx context.Context, thresholdSeconds int) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
