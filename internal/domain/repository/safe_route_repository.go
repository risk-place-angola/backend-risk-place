package repository

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type RouteCalculationParams struct {
	OriginLat      float64
	OriginLon      float64
	DestinationLat float64
	DestinationLon float64
	MaxRoutes      int
}

type IncidentHeatmapParams struct {
	NorthEastLat float64
	NorthEastLon float64
	SouthWestLat float64
	SouthWestLon float64
	StartDate    string
	EndDate      string
	RiskTypeID   string
}

type HeatmapPoint struct {
	Latitude     float64
	Longitude    float64
	Weight       float64
	IncidentType string
	ReportCount  int
}

type SafeRouteRepository interface {
	CalculateSafeRoute(ctx context.Context, params RouteCalculationParams) (*model.SafeRoute, error)
	GetIncidentsForRoute(ctx context.Context, waypoints []model.Waypoint, corridorWidthKm float64) ([]model.IncidentNearRoute, error)
	GetIncidentsHeatmap(ctx context.Context, params IncidentHeatmapParams) ([]HeatmapPoint, error)
}
