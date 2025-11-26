package service

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type DangerZoneParams struct {
	CenterLat    float64
	CenterLon    float64
	RadiusMeters float64
	MinIncidents int
}

type DangerZoneService interface {
	CalculateDangerZones(ctx context.Context) error
	GetDangerZonesNearby(ctx context.Context, lat, lon, radiusMeters float64) ([]*model.DangerZone, error)
	IsInDangerZone(ctx context.Context, lat, lon float64) (*model.DangerZone, error)
	InvalidateCache(ctx context.Context) error
}
