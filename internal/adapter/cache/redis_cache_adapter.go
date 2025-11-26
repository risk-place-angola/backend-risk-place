package cache

import (
	"context"
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/redis"
)

type redisCacheAdapter struct {
	redis *redis.Redis
}

func NewRedisCacheAdapter(rdb *redis.Redis) service.CacheService {
	return &redisCacheAdapter{redis: rdb}
}

func (a *redisCacheAdapter) Get(ctx context.Context, key string) (string, error) {
	return a.redis.Get(ctx, key)
}

func (a *redisCacheAdapter) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return a.redis.Set(ctx, key, value, ttl)
}

func (a *redisCacheAdapter) Delete(ctx context.Context, key string) error {
	return a.redis.Delete(ctx, key)
}

func (a *redisCacheAdapter) HGet(ctx context.Context, key, field string) (string, error) {
	return a.redis.HGet(ctx, key, field)
}

func (a *redisCacheAdapter) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return a.redis.HGetAll(ctx, key)
}

func (a *redisCacheAdapter) HSet(ctx context.Context, key, field, value string) error {
	return a.redis.HSet(ctx, key, field, value)
}

func (a *redisCacheAdapter) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) {
	a.redis.HExpire(ctx, key, expiration, fields...)
}

func (a *redisCacheAdapter) GeoAdd(ctx context.Context, key string, longitude, latitude float64, member string) error {
	return a.redis.GeoAdd(ctx, key, longitude, latitude, member)
}

func (a *redisCacheAdapter) GeoSearchWithDistance(ctx context.Context, key string, longitude, latitude, radiusMeters float64) ([]service.GeoResult, error) {
	portResults, err := a.redis.GeoSearchWithDistance(ctx, key, longitude, latitude, radiusMeters)
	if err != nil {
		return nil, err
	}

	results := make([]service.GeoResult, len(portResults))
	for i, r := range portResults {
		results[i] = service.GeoResult{
			Member:   r.Member,
			Distance: r.Distance,
		}
	}

	return results, nil
}
