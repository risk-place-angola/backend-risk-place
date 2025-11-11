package location

import (
	"context"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"log/slog"
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
