package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
)

type Redis struct {
	client *redis.Client
}

const (
	// Define any constants related to the database connection here
	RedisDefaultExpiration = 0 // 0 means the key has no expiration
	WithTimeout            = 50 * time.Second
)

// NewRedis creates a new instance of the Redis cache.
func NewRedis(cfg config.Config) *Redis {
	addr := fmt.Sprintf("%s:%d",
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Port,
	)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), WithTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	return &Redis{
		client: client,
	}
}

// Get returns the value of the key.
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}
	return res, nil
}

// HGet returns all fields and values of the hash stored at key.
func (r *Redis) HGet(ctx context.Context, key, field string) (string, error) {
	res, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

// HGetAll returns all fields and values of the hash stored at key.
func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	res, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Set stores a key-value pair in the Redis cache with an optional expiration timer.
// If the timer is less than or equal to zero, the key will persist indefinitely.
//
// Parameters:
//   - ctx: context.Context - The context for managing timeouts and cancellations.
//   - key: string - The key under which the value is stored.
//   - value: string - The value to store in Redis.
//   - timer: time.Duration - The expiration duration for the key. Use a value <= 0 for no expiration.
func (r *Redis) Set(ctx context.Context, key, value string, timer time.Duration) error {
	if timer <= RedisDefaultExpiration {
		timer = RedisDefaultExpiration
	}
	err := r.client.Set(ctx, key, value, timer).Err()
	return err
}

// HSet stores a key-value pair in the Redis cache with an optional expiration timer.
// If the timer is less than or equal to zero, the key will persist indefinitely.
//
// Parameters:
//   - ctx: context.Context - The context for managing timeouts and cancellations.
//   - key: string - The key under which the value is stored.
//   - field: string - The field under which the value is stored.
//   - value: string - The value to store in Redis.
func (r *Redis) HSet(ctx context.Context, key, field, value string) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// Delete removes the specified key from the Redis cache.
//
// Parameters:
//   - ctx: context.Context - The context for managing timeouts and cancellations.
//   - key: string - The key to remove from the cache.
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// HDelete removes the specified field from the hash stored at key.
//
// Parameters:
//   - ctx: context.Context - The context for managing timeouts and cancellations.
//   - key: string - The key under which the hash is stored.
//   - field: string - The field to remove from the hash.
func (r *Redis) HDelete(ctx context.Context, key, field string) error {
	return r.client.HDel(ctx, key, field).Err()
}

// Ping checks the health of the Redis connection by pinging the server.
// It returns true if the connection is healthy, and false otherwise.
func (r *Redis) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

// HExpire sets an expiration time for a field in a hash stored at key.
func (r *Redis) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) {
	r.client.HExpire(ctx, key, expiration, fields...)
}

// GeoAdd adds a geospatial item (member) with specified longitude and latitude to the geospatial index stored at key.
func (r *Redis) GeoAdd(ctx context.Context, key string, longitude, latitude float64, member string) error {
	return r.client.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      member,
		Longitude: longitude,
		Latitude:  latitude,
	}).Err()
}

// GeoSearch performs a geospatial search to find members within a specified radius from given longitude and latitude.
func (r *Redis) GeoSearch(ctx context.Context, key string, longitude float64, latitude float64, radiusMeters float64) ([]string, error) {
	res, err := r.client.GeoSearch(ctx, key, &redis.GeoSearchQuery{
		Longitude:  longitude,
		Latitude:   latitude,
		Radius:     radiusMeters,
		RadiusUnit: "m",
		Sort:       "ASC",
	}).Result()
	if err != nil {
		return nil, err
	}

	members := make([]string, len(res))
	for i, loc := range res {
		slog.Info("found member in radius", "member", loc)
		members[i] = loc
	}
	return members, nil
}

// GeoSearchWithDistance performs a geospatial search returning members with their distances.
// Results are already sorted by distance (ASC) by Redis.
func (r *Redis) GeoSearchWithDistance(ctx context.Context, key string, longitude float64, latitude float64, radiusMeters float64) ([]port.GeoResult, error) {
	// Usando GeoRadius com WITHDIST para obter distâncias calculadas pelo Redis
	res, err := r.client.GeoRadius(ctx, key, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:    radiusMeters,
		Unit:      "m",
		WithDist:  true,  // Retorna distâncias
		Sort:      "ASC", // Já ordena por distância
		Store:     "",
		StoreDist: "",
	}).Result()
	if err != nil {
		return nil, err
	}

	results := make([]port.GeoResult, len(res))
	for i, loc := range res {
		results[i] = port.GeoResult{
			Member:   loc.Name,
			Distance: loc.Dist,
		}
	}

	return results, nil
}

// GeoRemove removes a geospatial member from the geospatial index stored at key.
func (r *Redis) GeoRemove(ctx context.Context, key string, member string) error {
	return r.client.ZRem(ctx, key, member).Err()
}
