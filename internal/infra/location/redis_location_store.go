package location

import (
	"context"
	"log/slog"

	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
)

type RedisLocationStore struct {
	cache port.Cache
}

func NewRedisLocationStore(cache port.Cache) *RedisLocationStore {
	return &RedisLocationStore{
		cache: cache,
	}
}

func (s *RedisLocationStore) UpdateUserLocation(ctx context.Context, userID string, lat, lon float64) error {
	return s.cache.GeoAdd(ctx, "user_locations", lon, lat, userID)
}

func (s *RedisLocationStore) FindUsersInRadius(ctx context.Context, lat, lon float64, radiusMeters float64) ([]string, error) {
	users, err := s.cache.GeoSearch(ctx, "user_locations", lon, lat, radiusMeters)
	if err != nil {
		slog.Error("failed to find users in radius", "error", err)
		return nil, err
	}

	slog.Info("found users in radius", "count", len(users))

	return users, nil
}

// UpdateReportLocation adiciona ou atualiza a localização de um report no Redis
func (s *RedisLocationStore) UpdateReportLocation(ctx context.Context, reportID string, lat, lon float64) error {
	return s.cache.GeoAdd(ctx, "report_locations", lon, lat, reportID)
}

// FindReportsInRadius busca IDs de reports dentro de um raio específico
func (s *RedisLocationStore) FindReportsInRadius(ctx context.Context, lat, lon float64, radiusMeters float64) ([]string, error) {
	reports, err := s.cache.GeoSearch(ctx, "report_locations", lon, lat, radiusMeters)
	if err != nil {
		slog.Error("failed to find reports in radius", "error", err)
		return nil, err
	}

	slog.Info("found reports in radius", "count", len(reports))

	return reports, nil
}

// RemoveReportLocation remove a localização de um report do Redis (quando deletado)
func (s *RedisLocationStore) RemoveReportLocation(ctx context.Context, reportID string) error {
	return s.cache.GeoRemove(ctx, "report_locations", reportID)
}

// FindReportsInRadiusWithDistance busca reports com distâncias já calculadas e ordenadas pelo Redis
// Esta é a versão OTIMIZADA que evita cálculos redundantes
func (s *RedisLocationStore) FindReportsInRadiusWithDistance(ctx context.Context, lat, lon float64, radiusMeters float64) ([]port.GeoResult, error) {
	results, err := s.cache.GeoSearchWithDistance(ctx, "report_locations", lon, lat, radiusMeters)
	if err != nil {
		slog.Error("failed to find reports with distance", "error", err)
		return nil, err
	}

	// Converte o resultado do Redis para o tipo do port
	geoResults := make([]port.GeoResult, len(results))
	for i, r := range results {
		geoResults[i] = port.GeoResult{
			Member:   r.Member,
			Distance: r.Distance,
		}
	}

	slog.Info("found reports in radius with distances", "count", len(geoResults))
	return geoResults, nil
}
