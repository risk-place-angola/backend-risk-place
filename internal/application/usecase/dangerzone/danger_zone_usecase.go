package dangerzone

import (
	"context"
	"fmt"

	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type DangerZoneUseCase struct {
	dangerZoneService service.DangerZoneService
}

func NewDangerZoneUseCase(dangerZoneService service.DangerZoneService) *DangerZoneUseCase {
	return &DangerZoneUseCase{
		dangerZoneService: dangerZoneService,
	}
}

func (uc *DangerZoneUseCase) GetDangerZonesNearby(ctx context.Context, req *dto.GetDangerZonesRequest) (*dto.GetDangerZonesResponse, error) {
	zones, err := uc.dangerZoneService.GetDangerZonesNearby(ctx, req.Latitude, req.Longitude, req.RadiusMeters)
	if err != nil {
		return nil, fmt.Errorf("failed to get danger zones: %w", err)
	}

	response := &dto.GetDangerZonesResponse{
		Zones:      make([]dto.DangerZoneDTO, 0, len(zones)),
		TotalCount: len(zones),
	}

	for _, zone := range zones {
		response.Zones = append(response.Zones, dto.DangerZoneDTO{
			ID:           zone.ID.String(),
			Latitude:     zone.CellLat,
			Longitude:    zone.CellLon,
			GridCellID:   zone.GridCellID,
			IncidentCount: zone.IncidentCount,
			RiskScore:    zone.RiskScore,
			RiskLevel:    zone.RiskLevel,
			CalculatedAt: zone.CalculatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}
