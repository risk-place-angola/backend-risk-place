package locationsharing

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type LocationSharingUseCase struct {
	repo                 repository.LocationSharingRepository
	userRepo             repository.UserRepository
	anonymousSessionRepo repository.AnonymousSessionRepository
	geoService           port.GeolocationService
	config               *config.Config
}

func NewLocationSharingUseCase(
	repo repository.LocationSharingRepository,
	userRepo repository.UserRepository,
	anonymousSessionRepo repository.AnonymousSessionRepository,
	geoService port.GeolocationService,
	config *config.Config,
) *LocationSharingUseCase {
	return &LocationSharingUseCase{
		repo:                 repo,
		userRepo:             userRepo,
		anonymousSessionRepo: anonymousSessionRepo,
		geoService:           geoService,
		config:               config,
	}
}

func (uc *LocationSharingUseCase) CreateLocationSharingForUser(
	ctx context.Context,
	userID uuid.UUID,
	req dto.CreateLocationSharingRequest,
) (*dto.LocationSharingResponse, error) {
	if err := uc.geoService.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
		slog.Error("invalid coordinates for location sharing", "error", err)
		return nil, err
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get user for location sharing", "error", err)
		return nil, errors.New("user not found")
	}

	sharing := model.NewLocationSharing(req.Latitude, req.Longitude, req.DurationMinutes, user.Name)
	sharing.SetAuthenticatedUser(userID)

	if err := uc.repo.Save(ctx, sharing); err != nil {
		slog.Error("failed to create location sharing", "error", err)
		return nil, err
	}

	baseURL := uc.config.FrontendURL
	if baseURL == "" {
		baseURL = "https://riskplace.ao"
	}

	slog.Info("location sharing created successfully", "sharing_id", sharing.ID, "user_id", userID)
	return dto.ToLocationSharingResponse(sharing, baseURL), nil
}

func (uc *LocationSharingUseCase) CreateLocationSharingForAnonymous(
	ctx context.Context,
	deviceID string,
	req dto.CreateLocationSharingRequest,
) (*dto.LocationSharingResponse, error) {
	if err := uc.geoService.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
		slog.Error("invalid coordinates for location sharing", "error", err)
		return nil, err
	}

	session, err := uc.anonymousSessionRepo.FindByDeviceID(ctx, deviceID)
	if err != nil {
		slog.Error("failed to get anonymous session", "error", err, "device_id", deviceID)
		return nil, errors.New("anonymous session not found")
	}

	ownerName := "Usuário Anônimo"
	sharing := model.NewLocationSharing(req.Latitude, req.Longitude, req.DurationMinutes, ownerName)
	sharing.SetAnonymousUser(session.ID, deviceID)

	if err := uc.repo.Save(ctx, sharing); err != nil {
		slog.Error("failed to create location sharing", "error", err)
		return nil, err
	}

	baseURL := uc.config.FrontendURL
	if baseURL == "" {
		baseURL = "https://riskplace.ao"
	}

	slog.Info("location sharing created successfully for anonymous", "sharing_id", sharing.ID, "device_id", deviceID)
	return dto.ToLocationSharingResponse(sharing, baseURL), nil
}

func (uc *LocationSharingUseCase) GetLocationSharingByToken(
	ctx context.Context,
	token string,
) (*dto.PublicLocationResponse, error) {
	sharing, err := uc.repo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("location sharing not found", "token", token)
			return nil, errors.New("location sharing not found")
		}
		slog.Error("failed to get location sharing", "error", err)
		return nil, err
	}

	if !sharing.IsValid() {
		slog.Warn("location sharing is expired or inactive", "token", token)
		return nil, errors.New("location sharing is expired or inactive")
	}

	return dto.ToPublicLocationResponse(sharing), nil
}

func (uc *LocationSharingUseCase) UpdateLocationForUser(
	ctx context.Context,
	sharingID uuid.UUID,
	userID uuid.UUID,
	req dto.UpdateLocationRequest,
) error {
	if err := uc.geoService.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
		slog.Error("invalid coordinates for location update", "error", err)
		return err
	}

	sharing, err := uc.repo.FindByID(ctx, sharingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("location sharing not found", "sharing_id", sharingID)
			return errors.New("location sharing not found")
		}
		slog.Error("failed to get location sharing", "error", err)
		return err
	}

	if !sharing.IsOwnedByUser(userID) {
		slog.Warn("unauthorized location update attempt", "sharing_id", sharingID, "user_id", userID)
		return errors.New("unauthorized")
	}

	if !sharing.IsValid() {
		slog.Warn("cannot update expired or inactive location sharing", "sharing_id", sharingID)
		return errors.New("location sharing is expired or inactive")
	}

	sharing.UpdateLocation(req.Latitude, req.Longitude)

	if err := uc.repo.Update(ctx, sharing); err != nil {
		slog.Error("failed to update location", "error", err)
		return err
	}

	slog.Info("location updated successfully", "sharing_id", sharingID)
	return nil
}

func (uc *LocationSharingUseCase) UpdateLocationForAnonymous(
	ctx context.Context,
	sharingID uuid.UUID,
	deviceID string,
	req dto.UpdateLocationRequest,
) error {
	if err := uc.geoService.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
		slog.Error("invalid coordinates for location update", "error", err)
		return err
	}

	sharing, err := uc.repo.FindByID(ctx, sharingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("location sharing not found", "sharing_id", sharingID)
			return errors.New("location sharing not found")
		}
		slog.Error("failed to get location sharing", "error", err)
		return err
	}

	if !sharing.IsOwnedByDevice(deviceID) {
		slog.Warn("unauthorized location update attempt", "sharing_id", sharingID, "device_id", deviceID)
		return errors.New("unauthorized")
	}

	if !sharing.IsValid() {
		slog.Warn("cannot update expired or inactive location sharing", "sharing_id", sharingID)
		return errors.New("location sharing is expired or inactive")
	}

	sharing.UpdateLocation(req.Latitude, req.Longitude)

	if err := uc.repo.Update(ctx, sharing); err != nil {
		slog.Error("failed to update location", "error", err)
		return err
	}

	slog.Info("location updated successfully for anonymous", "sharing_id", sharingID, "device_id", deviceID)
	return nil
}

func (uc *LocationSharingUseCase) DeleteLocationSharingForUser(
	ctx context.Context,
	sharingID uuid.UUID,
	userID uuid.UUID,
) error {
	sharing, err := uc.repo.FindByID(ctx, sharingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("location sharing not found", "sharing_id", sharingID)
			return errors.New("location sharing not found")
		}
		slog.Error("failed to get location sharing", "error", err)
		return err
	}

	if !sharing.IsOwnedByUser(userID) {
		slog.Warn("unauthorized delete attempt", "sharing_id", sharingID, "user_id", userID)
		return errors.New("unauthorized")
	}

	sharing.Deactivate()

	if err := uc.repo.Update(ctx, sharing); err != nil {
		slog.Error("failed to deactivate location sharing", "error", err)
		return err
	}

	slog.Info("location sharing deactivated successfully", "sharing_id", sharingID)
	return nil
}

func (uc *LocationSharingUseCase) DeleteLocationSharingForAnonymous(
	ctx context.Context,
	sharingID uuid.UUID,
	deviceID string,
) error {
	sharing, err := uc.repo.FindByID(ctx, sharingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("location sharing not found", "sharing_id", sharingID)
			return errors.New("location sharing not found")
		}
		slog.Error("failed to get location sharing", "error", err)
		return err
	}

	if !sharing.IsOwnedByDevice(deviceID) {
		slog.Warn("unauthorized delete attempt", "sharing_id", sharingID, "device_id", deviceID)
		return errors.New("unauthorized")
	}

	sharing.Deactivate()

	if err := uc.repo.Update(ctx, sharing); err != nil {
		slog.Error("failed to deactivate location sharing", "error", err)
		return err
	}

	slog.Info("location sharing deactivated successfully for anonymous", "sharing_id", sharingID, "device_id", deviceID)
	return nil
}
