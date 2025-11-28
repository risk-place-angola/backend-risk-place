package repository

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type DangerZoneCalculationResult struct {
	GridCellID   string
	CellLat      float64
	CellLon      float64
	IncidentCount int
	RiskScore    float64
}

type DangerZoneRepository interface {
	CalculateDangerZones(ctx context.Context, gridSizeKm float64, minIncidents int, daysBack int) ([]DangerZoneCalculationResult, error)
	GetNearbyDangerZones(ctx context.Context, lat, lon, radiusKm float64) ([]*model.DangerZone, error)
}
