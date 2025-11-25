package service

import (
	"context"
	"fmt"
	"log/slog"

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
		slog.Error("[ADAPTER] ❌ FAILED TO PARSE USER_ID AS UUID",
			slog.String("user_id", userID),
			slog.String("device_id", deviceID),
			slog.Bool("is_anonymous", isAnonymous),
			slog.Any("error", err))
		return fmt.Errorf("invalid user_id format: %w", err)
	}

	return a.domainService.UpdateUserLocation(ctx, uid, deviceID, lat, lon, speed, heading, isAnonymous)
}

func (a *nearbyUsersAdapter) GetNearbyUsers(ctx context.Context, userID string, lat, lon, radiusMeters float64) ([]port.NearbyUser, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("[ADAPTER] ❌ FAILED TO PARSE USER_ID AS UUID in GetNearbyUsers",
			slog.String("user_id", userID),
			slog.Any("error", err))
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	users, err := a.domainService.GetNearbyUsers(ctx, uid, lat, lon, radiusMeters)
	if err != nil {
		slog.Error("[ADAPTER] failed to get nearby users from domain service",
			slog.String("user_id", uid.String()),
			slog.Any("error", err))
		return nil, err
	}

	slog.Debug("[ADAPTER] Converting domain users to port users",
		slog.String("requesting_user", uid.String()),
		slog.Int("users_found", len(users)))

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

		slog.Debug("[ADAPTER] Converted user",
			slog.String("anonymous_id", u.AnonymousID),
			slog.Float64("lat", u.Latitude),
			slog.Float64("lon", u.Longitude))
	}

	return result, nil
}

func (a *nearbyUsersAdapter) CleanupStaleLocations(ctx context.Context) error {
	return a.domainService.CleanupStaleLocations(ctx)
}
