package saferoute

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type SafeRouteUseCase struct {
	safeRouteRepo repository.SafeRouteRepository
	userRepo      repository.UserRepository
}

func NewSafeRouteUseCase(safeRouteRepo repository.SafeRouteRepository, userRepo repository.UserRepository) *SafeRouteUseCase {
	return &SafeRouteUseCase{
		safeRouteRepo: safeRouteRepo,
		userRepo:      userRepo,
	}
}

func (uc *SafeRouteUseCase) CalculateSafeRoute(ctx context.Context, req *dto.SafeRouteRequest) (*dto.SafeRouteResponse, error) {
	if req.MaxRoutes == 0 {
		req.MaxRoutes = 1
	}

	params := repository.RouteCalculationParams{
		OriginLat:      req.OriginLat,
		OriginLon:      req.OriginLon,
		DestinationLat: req.DestinationLat,
		DestinationLon: req.DestinationLon,
		MaxRoutes:      req.MaxRoutes,
	}

	route, err := uc.safeRouteRepo.CalculateSafeRoute(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate safe route: %w", err)
	}

	return uc.toDTO(route), nil
}

func (uc *SafeRouteUseCase) GetIncidentsHeatmap(ctx context.Context, req *dto.HeatmapRequest) (*dto.HeatmapResponse, error) {
	params := repository.IncidentHeatmapParams{
		NorthEastLat: req.NorthEastLat,
		NorthEastLon: req.NorthEastLon,
		SouthWestLat: req.SouthWestLat,
		SouthWestLon: req.SouthWestLon,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		RiskTypeID:   req.RiskTypeID,
	}

	points, err := uc.safeRouteRepo.GetIncidentsHeatmap(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get incidents heatmap: %w", err)
	}

	response := &dto.HeatmapResponse{
		Points:     make([]dto.HeatmapPointDTO, 0, len(points)),
		TotalCount: len(points),
	}

	response.BoundsInfo.NorthEastLat = req.NorthEastLat
	response.BoundsInfo.NorthEastLon = req.NorthEastLon
	response.BoundsInfo.SouthWestLat = req.SouthWestLat
	response.BoundsInfo.SouthWestLon = req.SouthWestLon

	for _, point := range points {
		response.Points = append(response.Points, dto.HeatmapPointDTO{
			Latitude:     point.Latitude,
			Longitude:    point.Longitude,
			Weight:       point.Weight,
			IncidentType: point.IncidentType,
			ReportCount:  point.ReportCount,
		})
	}

	return response, nil
}

func (uc *SafeRouteUseCase) toDTO(route *model.SafeRoute) *dto.SafeRouteResponse {
	waypoints := make([]dto.WaypointDTO, 0, len(route.Waypoints))
	for _, wp := range route.Waypoints {
		waypoints = append(waypoints, dto.WaypointDTO{
			Latitude:  wp.Latitude,
			Longitude: wp.Longitude,
			Sequence:  wp.Sequence,
		})
	}

	incidents := make([]dto.IncidentDTO, 0, len(route.Incidents))
	for _, inc := range route.Incidents {
		incidents = append(incidents, dto.IncidentDTO{
			ReportID:     inc.ReportID.String(),
			RiskType:     inc.RiskType,
			RiskTopic:    inc.RiskTopic,
			Latitude:     inc.Latitude,
			Longitude:    inc.Longitude,
			DistanceKm:   inc.DistanceKm,
			CreatedAt:    inc.CreatedAt,
			DaysAgo:      inc.DaysAgo,
			WeightFactor: inc.WeightFactor,
		})
	}

	return &dto.SafeRouteResponse{
		ID:                route.ID.String(),
		OriginLat:         route.OriginLat,
		OriginLon:         route.OriginLon,
		DestinationLat:    route.DestinationLat,
		DestinationLon:    route.DestinationLon,
		Waypoints:         waypoints,
		DistanceKm:        route.DistanceKm,
		EstimatedDuration: route.EstimatedDuration,
		SafetyScore:       route.SafetyScore,
		RiskLevel:         string(route.RiskLevel),
		IncidentCount:     route.IncidentCount,
		Incidents:         incidents,
		CalculatedAt:      route.CalculatedAt,
	}
}

// navigateToSavedLocation is a helper to reduce duplication between NavigateToHome and NavigateToWork
func (uc *SafeRouteUseCase) navigateToSavedLocation(
	ctx context.Context,
	userID uuid.UUID,
	currentLat, currentLon float64,
	getAddress func(user *model.User) *model.SavedLocation,
	notConfiguredMsg string,
) (*dto.SafeRouteResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	address := getAddress(user)
	if address == nil {
		return nil, errors.New(notConfiguredMsg)
	}

	req := &dto.SafeRouteRequest{
		OriginLat:      currentLat,
		OriginLon:      currentLon,
		DestinationLat: address.Latitude,
		DestinationLon: address.Longitude,
		MaxRoutes:      1,
	}

	return uc.CalculateSafeRoute(ctx, req)
}

func (uc *SafeRouteUseCase) NavigateToHome(ctx context.Context, userID uuid.UUID, currentLat, currentLon float64) (*dto.SafeRouteResponse, error) {
	return uc.navigateToSavedLocation(ctx, userID, currentLat, currentLon,
		func(user *model.User) *model.SavedLocation { return user.HomeAddress },
		"home address not configured",
	)
}

func (uc *SafeRouteUseCase) NavigateToWork(ctx context.Context, userID uuid.UUID, currentLat, currentLon float64) (*dto.SafeRouteResponse, error) {
	return uc.navigateToSavedLocation(ctx, userID, currentLat, currentLon,
		func(user *model.User) *model.SavedLocation { return user.WorkAddress },
		"work address not configured",
	)
}
