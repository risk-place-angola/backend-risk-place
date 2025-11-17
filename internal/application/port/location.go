package port

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type LocationStore interface {
	UpdateUserLocation(ctx context.Context, userID string, lat float64, lon float64) error
	FindUsersInRadius(ctx context.Context, lat float64, lon float64, radiusMeters float64) ([]string, error)
	RemoveReportLocation(ctx context.Context, reportID string) error
	UpdateReportLocation(ctx context.Context, reportID string, lat, lon float64) error
	FindReportsInRadius(ctx context.Context, lat, lon float64, radiusMeters float64) ([]string, error)
}

type GeolocationService interface {
	ValidateCoordinates(lat, lon float64) error
	DistanceBetween(p1, p2 service.Geolocation) float64
	IsWithinRadius(p1, p2 service.Geolocation, radiusMeters float64) bool
	Parse(lat, lon float64) (service.Geolocation, error)
	Midpoint(p1, p2 service.Geolocation) service.Geolocation
}
