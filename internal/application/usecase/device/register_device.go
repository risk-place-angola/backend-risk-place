package device

import (
	"context"
	"fmt"

	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type RegisterDeviceUseCase struct {
	anonymousSessionRepo repository.AnonymousSessionRepository
}

func NewRegisterDeviceUseCase(repo repository.AnonymousSessionRepository) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{
		anonymousSessionRepo: repo,
	}
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, req dto.RegisterDeviceRequest) (*dto.DeviceResponse, error) {
	existingSession, err := uc.anonymousSessionRepo.FindByDeviceID(ctx, req.DeviceID)

	if err == nil && existingSession != nil {
		return uc.updateExistingSession(ctx, existingSession, req)
	}

	session, err := model.NewAnonymousSession(
		req.DeviceID,
		req.FCMToken,
		req.Platform,
		req.Model,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create anonymous session: %w", err)
	}

	if req.Language != "" {
		session.DeviceLanguage = req.Language
	}

	if req.Latitude != 0 && req.Longitude != 0 {
		session.UpdateLocation(req.Latitude, req.Longitude)
	}

	if req.AlertRadiusMeters > 0 {
		session.AlertRadiusMeters = req.AlertRadiusMeters
	}

	err = uc.anonymousSessionRepo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create device session: %w", err)
	}

	return &dto.DeviceResponse{
		DeviceID:          session.DeviceID,
		FCMToken:          session.DeviceFCMToken,
		Platform:          session.DevicePlatform,
		Latitude:          session.Latitude,
		Longitude:         session.Longitude,
		AlertRadiusMeters: session.AlertRadiusMeters,
		Message:           "Device registered successfully",
	}, nil
}

func (uc *RegisterDeviceUseCase) updateExistingSession(ctx context.Context, session *model.AnonymousSession, req dto.RegisterDeviceRequest) (*dto.DeviceResponse, error) {
	if req.FCMToken != "" {
		session.UpdateFCMToken(req.FCMToken)
	}
	if req.Latitude != 0 && req.Longitude != 0 {
		session.UpdateLocation(req.Latitude, req.Longitude)
	}
	if req.Platform != "" {
		session.DevicePlatform = req.Platform
	}
	if req.Model != "" {
		session.DeviceModel = req.Model
	}
	if req.Language != "" {
		session.DeviceLanguage = req.Language
	}
	if req.AlertRadiusMeters > 0 {
		session.AlertRadiusMeters = req.AlertRadiusMeters
	}

	err := uc.anonymousSessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update device session: %w", err)
	}

	return &dto.DeviceResponse{
		DeviceID:          session.DeviceID,
		FCMToken:          session.DeviceFCMToken,
		Platform:          session.DevicePlatform,
		Latitude:          session.Latitude,
		Longitude:         session.Longitude,
		AlertRadiusMeters: session.AlertRadiusMeters,
		Message:           "Device session updated successfully",
	}, nil
}
