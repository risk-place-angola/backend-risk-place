package service

import (
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type GeolocationAdapter struct {
	domainService *service.DefaultGeolocationService
}

func NewGeolocationAdapter(domainService *service.DefaultGeolocationService) port.GeolocationService {
	return &GeolocationAdapter{
		domainService: domainService,
	}
}

func (a *GeolocationAdapter) ValidateCoordinates(lat, lon float64) error {
	return a.domainService.ValidateCoordinates(lat, lon)
}

func (a *GeolocationAdapter) DistanceBetween(p1, p2 port.Geolocation) float64 {
	domainP1 := service.Geolocation{Latitude: p1.Latitude, Longitude: p1.Longitude}
	domainP2 := service.Geolocation{Latitude: p2.Latitude, Longitude: p2.Longitude}
	return a.domainService.DistanceBetween(domainP1, domainP2)
}

func (a *GeolocationAdapter) IsWithinRadius(p1, p2 port.Geolocation, radiusMeters float64) bool {
	domainP1 := service.Geolocation{Latitude: p1.Latitude, Longitude: p1.Longitude}
	domainP2 := service.Geolocation{Latitude: p2.Latitude, Longitude: p2.Longitude}
	return a.domainService.IsWithinRadius(domainP1, domainP2, radiusMeters)
}

func (a *GeolocationAdapter) Parse(lat, lon float64) (port.Geolocation, error) {
	domainGeo, err := a.domainService.Parse(lat, lon)
	if err != nil {
		return port.Geolocation{}, err
	}
	return port.Geolocation{Latitude: domainGeo.Latitude, Longitude: domainGeo.Longitude}, nil
}

func (a *GeolocationAdapter) Midpoint(p1, p2 port.Geolocation) port.Geolocation {
	domainP1 := service.Geolocation{Latitude: p1.Latitude, Longitude: p1.Longitude}
	domainP2 := service.Geolocation{Latitude: p2.Latitude, Longitude: p2.Longitude}
	domainMidpoint := a.domainService.Midpoint(domainP1, domainP2)
	return port.Geolocation{Latitude: domainMidpoint.Latitude, Longitude: domainMidpoint.Longitude}
}
