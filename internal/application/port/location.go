package port

import (
	"context"
)

type GeoResult struct {
	Member   string
	Distance float64
}

type Geolocation struct {
	Latitude  float64
	Longitude float64
}

type LocationStore interface {
	UpdateUserLocation(ctx context.Context, userID string, lat float64, lon float64) error
	FindUsersInRadius(ctx context.Context, lat float64, lon float64, radiusMeters float64) ([]string, error)
	RemoveReportLocation(ctx context.Context, reportID string) error
	UpdateReportLocation(ctx context.Context, reportID string, lat, lon float64) error
	FindReportsInRadius(ctx context.Context, lat, lon float64, radiusMeters float64) ([]string, error)
	FindReportsInRadiusWithDistance(ctx context.Context, lat, lon float64, radiusMeters float64) ([]GeoResult, error)
}

type GeolocationService interface {
	ValidateCoordinates(lat, lon float64) error
	DistanceBetween(p1, p2 Geolocation) float64
	IsWithinRadius(p1, p2 Geolocation, radiusMeters float64) bool
	Parse(lat, lon float64) (Geolocation, error)
	Midpoint(p1, p2 Geolocation) Geolocation
}
