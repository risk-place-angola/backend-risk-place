package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type nearbyUsersAdapter struct {
	domainService service.NearbyUsersService
}

func NewNearbyUsersAdapter(domainService service.NearbyUsersService) port.NearbyUsersService {
	return &nearbyUsersAdapter{domainService: domainService}
}

func (a *nearbyUsersAdapter) UpdateUserLocation(ctx context.Context, userID string, deviceID string, lat, lon, speed, heading float64, isAnonymous bool) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.domainService.UpdateUserLocation(ctx, uid, deviceID, lat, lon, speed, heading, isAnonymous)
}

func (a *nearbyUsersAdapter) GetNearbyUsers(ctx context.Context, userID string, lat, lon, radiusMeters float64) ([]port.NearbyUser, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	users, err := a.domainService.GetNearbyUsers(ctx, uid, lat, lon, radiusMeters)
	if err != nil {
		return nil, err
	}

	result := make([]port.NearbyUser, len(users))
	for i, u := range users {
		result[i] = port.NearbyUser{
			UserID:      u.AnonymousID,
			AnonymousID: u.AnonymousID,
			Latitude:    u.Latitude,
			Longitude:   u.Longitude,
			AvatarID:    u.AvatarID,
			Color:       u.Color,
			Speed:       u.Speed,
			Heading:     u.Heading,
		}
	}

	return result, nil
}

func (a *nearbyUsersAdapter) CleanupStaleLocations(ctx context.Context) error {
	return a.domainService.CleanupStaleLocations(ctx)
}
