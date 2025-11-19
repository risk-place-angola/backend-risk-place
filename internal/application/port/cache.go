package port

import (
	"context"
	"time"
)

type Cache interface {
	KVCache
	HashCache
	GeoCache
	HealthChecker
}

type KVCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, timer time.Duration) error
	Delete(ctx context.Context, key string) error
}

type HashCache interface {
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key, field, value string) error
	HDelete(ctx context.Context, key, field string) error
	HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string)
}

type GeoCache interface {
	GeoAdd(ctx context.Context, key string, longitude, latitude float64, member string) error
	GeoSearch(ctx context.Context, key string, longitude, latitude, radiusMeters float64) ([]string, error)
	GeoSearchWithDistance(ctx context.Context, key string, longitude, latitude, radiusMeters float64) ([]GeoResult, error)
	GeoRemove(ctx context.Context, key string, member string) error
}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type InMemory interface {
	Get(key string, value any) error
	Set(key string, value any) error
	SetWithTTL(key string, value any, ttl int) error
	Delete(key string) error
	Reset() error
	Metrics()
}
