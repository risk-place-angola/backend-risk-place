package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

const (
	staleLocationThresholdSeconds = 30
)

type NearbyUsersService interface {
	UpdateUserLocation(ctx context.Context, userID uuid.UUID, deviceID string, lat, lon, speed, heading float64, isAnonymous bool) error
	GetNearbyUsers(ctx context.Context, userID uuid.UUID, lat, lon, radiusMeters float64) ([]*model.NearbyUser, error)
	CleanupStaleLocations(ctx context.Context) error
}

type nearbyUsersService struct {
	repo repository.UserLocationRepository
}

func NewNearbyUsersService(repo repository.UserLocationRepository) NearbyUsersService {
	return &nearbyUsersService{repo: repo}
}

func (s *nearbyUsersService) UpdateUserLocation(ctx context.Context, userID uuid.UUID, deviceID string, lat, lon, speed, heading float64, isAnonymous bool) error {
	location := model.NewUserLocation(userID, deviceID, lat, lon, speed, heading, isAnonymous)

	if err := s.repo.Upsert(ctx, location); err != nil {
		return err
	}

	// Note: Location history now handled by Redis-based LocationHistoryService
	// See: internal/adapter/service/location_history_service.go

	return nil
}

func (s *nearbyUsersService) GetNearbyUsers(ctx context.Context, requestingUserID uuid.UUID, lat, lon, radiusMeters float64) ([]*model.NearbyUser, error) {
	const maxUsers = 100

	locations, err := s.repo.FindNearbyUsers(ctx, lat, lon, radiusMeters, maxUsers+1)
	if err != nil {
		return nil, err
	}

	nearbyUsers := make([]*model.NearbyUser, 0, len(locations))
	for _, loc := range locations {
		if loc.UserID == requestingUserID {
			continue
		}

		privacyLat, privacyLon := model.ApplyPrivacyOffset(loc.Latitude, loc.Longitude)

		avatarID := fmt.Sprintf("avatar_%d", loc.AvatarID)

		nearbyUser := &model.NearbyUser{
			UserID:      loc.UserID,
			AnonymousID: model.GenerateAnonymousID(loc.UserID),
			Latitude:    privacyLat,
			Longitude:   privacyLon,
			AvatarID:    avatarID,
			Color:       loc.Color,
			Speed:       loc.Speed,
			Heading:     loc.Heading,
			LastUpdate:  loc.LastUpdate,
			IsAnonymous: loc.IsAnonymous,
		}

		nearbyUsers = append(nearbyUsers, nearbyUser)

		if len(nearbyUsers) >= maxUsers {
			break
		}
	}

	return nearbyUsers, nil
}

func (s *nearbyUsersService) CleanupStaleLocations(ctx context.Context) error {
	return s.repo.DeleteStale(ctx, staleLocationThresholdSeconds)
}
