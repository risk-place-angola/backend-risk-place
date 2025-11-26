package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

const (
	dangerZoneCacheKey       = "danger_zones"
	dangerZoneCacheTTL       = 1 * time.Hour
	dangerZoneGridSize       = 0.5
	dangerZoneMinIncidents   = 5
	dangerZoneDaysBack       = 30
	dangerZoneRecalcInterval = 30 * time.Minute
	dangerZoneRadiusMeters   = 500.0
	metersToKm               = 1000.0
)

type DangerZoneServiceImpl struct {
	repo  repository.DangerZoneRepository
	cache domainService.CacheService
}

func NewDangerZoneService(
	repo repository.DangerZoneRepository,
	cache domainService.CacheService,
) domainService.DangerZoneService {
	return &DangerZoneServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

func (s *DangerZoneServiceImpl) CalculateDangerZones(ctx context.Context) error {
	results, err := s.repo.CalculateDangerZones(ctx, dangerZoneGridSize, dangerZoneMinIncidents, dangerZoneDaysBack)
	if err != nil {
		return fmt.Errorf("failed to calculate danger zones: %w", err)
	}

	for _, result := range results {
		zone := model.NewDangerZone(result.CellLat, result.CellLon, result.GridCellID)
		zone.IncidentCount = result.IncidentCount
		zone.RiskScore = result.RiskScore
		zone.CalculateRiskLevel()

		data, err := json.Marshal(zone)
		if err != nil {
			slog.Debug("failed to marshal danger zone", "error", err, "grid_cell_id", result.GridCellID)
			continue
		}

		key := fmt.Sprintf("%s:%s", dangerZoneCacheKey, result.GridCellID)
		if err := s.cache.Set(ctx, key, string(data), dangerZoneCacheTTL); err != nil {
			slog.Debug("failed to cache danger zone", "error", err, "grid_cell_id", result.GridCellID)
		}

		if err := s.cache.GeoAdd(ctx, dangerZoneCacheKey, result.CellLon, result.CellLat, result.GridCellID); err != nil {
			slog.Debug("failed to add danger zone to geospatial index", "error", err, "grid_cell_id", result.GridCellID)
		}
	}

	return nil
}

func (s *DangerZoneServiceImpl) GetDangerZonesNearby(ctx context.Context, lat, lon, radiusMeters float64) ([]*model.DangerZone, error) {
	geoResults, err := s.cache.GeoSearchWithDistance(ctx, dangerZoneCacheKey, lon, lat, radiusMeters)
	if err != nil {
		slog.Debug("cache miss for danger zones, querying database", "error", err)
		return s.repo.GetNearbyDangerZones(ctx, lat, lon, radiusMeters/metersToKm)
	}

	if len(geoResults) == 0 {
		return s.repo.GetNearbyDangerZones(ctx, lat, lon, radiusMeters/metersToKm)
	}

	var zones []*model.DangerZone
	for _, result := range geoResults {
		key := fmt.Sprintf("%s:%s", dangerZoneCacheKey, result.Member)
		data, err := s.cache.Get(ctx, key)
		if err != nil {
			continue
		}

		var zone model.DangerZone
		if err := json.Unmarshal([]byte(data), &zone); err != nil {
			continue
		}

		if !zone.IsExpired() {
			zones = append(zones, &zone)
		}
	}

	return zones, nil
}

var ErrNoDangerZoneFound = fmt.Errorf("no danger zone found")

func (s *DangerZoneServiceImpl) IsInDangerZone(ctx context.Context, lat, lon float64) (*model.DangerZone, error) {
	zones, err := s.GetDangerZonesNearby(ctx, lat, lon, dangerZoneRadiusMeters)
	if err != nil {
		return nil, err
	}

	for _, zone := range zones {
		if zone.RiskLevel == "high" || zone.RiskLevel == "critical" {
			return zone, nil
		}
	}

	return nil, ErrNoDangerZoneFound
}

func (s *DangerZoneServiceImpl) InvalidateCache(ctx context.Context) error {
	if err := s.cache.Delete(ctx, dangerZoneCacheKey); err != nil {
		return fmt.Errorf("failed to invalidate danger zone cache: %w", err)
	}
	return nil
}
