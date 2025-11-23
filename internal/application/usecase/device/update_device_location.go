package device

import (
	"context"
	"fmt"

	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type UpdateDeviceLocationUseCase struct {
	anonymousSessionRepo repository.AnonymousSessionRepository
	locationStore        port.LocationStore
}

func NewUpdateDeviceLocationUseCase(
	repo repository.AnonymousSessionRepository,
	locationStore port.LocationStore,
) *UpdateDeviceLocationUseCase {
	return &UpdateDeviceLocationUseCase{
		anonymousSessionRepo: repo,
		locationStore:        locationStore,
	}
}

func (uc *UpdateDeviceLocationUseCase) Execute(ctx context.Context, req dto.UpdateDeviceLocationRequest) error {
	err := uc.anonymousSessionRepo.UpdateLocation(ctx, req.DeviceID, req.Latitude, req.Longitude)
	if err != nil {
		return fmt.Errorf("failed to update device location: %w", err)
	}

	err = uc.locationStore.UpdateUserLocation(ctx, req.DeviceID, req.Latitude, req.Longitude)
	if err != nil {
		return fmt.Errorf("failed to update location store: %w", err)
	}

	return nil
}
