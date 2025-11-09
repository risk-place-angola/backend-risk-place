package service

import (
	"errors"
	"math"
)

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DefaultGeolocationService is the standard implementation of the service.
type DefaultGeolocationService struct{}

func NewGeolocationService() *DefaultGeolocationService {
	return &DefaultGeolocationService{}
}

// ValidateCoordinates checks whether latitude and longitude are valid.
func (g *DefaultGeolocationService) ValidateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude value")
	}
	if lon < -180 || lon > 180 {
		return errors.New("invalid longitude value")
	}
	return nil
}

// DistanceBetween calculates the distance in meters between two geolocations using the Haversine formula.
func (g *DefaultGeolocationService) DistanceBetween(p1, p2 Geolocation) float64 {
	const EarthRadius = 6371e3

	lat1 := p1.Latitude * math.Pi / 180
	lat2 := p2.Latitude * math.Pi / 180
	dLat := (p2.Latitude - p1.Latitude) * math.Pi / 180
	dLon := (p2.Longitude - p1.Longitude) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * c
}

// IsWithinRadius checks if two geolocations are within a specified radius in meters.
func (g *DefaultGeolocationService) IsWithinRadius(p1, p2 Geolocation, radiusMeters float64) bool {
	return g.DistanceBetween(p1, p2) <= radiusMeters
}

// Parse creates a Geolocation from latitude and longitude, validating the inputs.
func (g *DefaultGeolocationService) Parse(lat, lon float64) (Geolocation, error) {
	if err := g.ValidateCoordinates(lat, lon); err != nil {
		return Geolocation{}, err
	}
	return Geolocation{Latitude: lat, Longitude: lon}, nil
}

// Midpoint calculates the midpoint between two geolocations
func (g *DefaultGeolocationService) Midpoint(p1, p2 Geolocation) Geolocation {
	lat1 := p1.Latitude * math.Pi / 180
	lon1 := p1.Longitude * math.Pi / 180
	lat2 := p2.Latitude * math.Pi / 180
	dLon := (p2.Longitude - p1.Longitude) * math.Pi / 180

	bx := math.Cos(lat2) * math.Cos(dLon)
	by := math.Cos(lat2) * math.Sin(dLon)

	lat3 := math.Atan2(
		math.Sin(lat1)+math.Sin(lat2),
		math.Sqrt((math.Cos(lat1)+bx)*(math.Cos(lat1)+bx)+by*by),
	)
	lon3 := lon1 + math.Atan2(by, math.Cos(lat1)+bx)

	return Geolocation{
		Latitude:  lat3 * 180 / math.Pi,
		Longitude: lon3 * 180 / math.Pi,
	}
}
