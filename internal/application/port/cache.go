package port

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Set(ctx context.Context, key, value string, timer time.Duration) error
	HSet(ctx context.Context, key, field, value string) error
	Delete(ctx context.Context, key string) error
	HDelete(ctx context.Context, key, field string) error
	Ping(ctx context.Context) error
	HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string)
	GeoAdd(ctx context.Context, key string, longitude, latitude float64, member string) error
	GeoSearch(ctx context.Context, key string, longitude float64, latitude float64, radiusMeters float64) ([]string, error)
}

type InMemory interface {
	Get(key string, value any) error
	Set(key string, value any) error
	SetWithTTL(key string, value any, ttl int) error
	Delete(key string) error
	Reset() error
	Metrics()
}
