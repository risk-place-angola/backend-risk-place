package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

// GeoResult represents a geo search result with distance
type GeoResult struct {
	Member   string
	Distance float64
}

type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key, field, value string) error
	HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string)
	GeoAdd(ctx context.Context, key string, longitude, latitude float64, member string) error
	GeoSearchWithDistance(ctx context.Context, key string, longitude, latitude, radiusMeters float64) ([]GeoResult, error)
}

const (
	nearbyUsersCacheTTLV2   = 5 * time.Second
	locationMetadataTTL     = 60 * time.Second
	redisGeoKey             = "user_locations:geo"
	maxConcurrentGoroutines = 50
	maxNearbyUsersLimit     = 100
	percentageMultiplier    = 100
)

type NearbyUsersServiceV2 struct {
	repo         repository.UserLocationRepository
	cache        CacheService
	useRedis     bool
	fallbackToPG bool
	cacheHits    int64
	cacheMisses  int64
	redisErrors  int64
}

func NewNearbyUsersServiceV2(
	repo repository.UserLocationRepository,
	cache CacheService,
	useRedis bool,
) NearbyUsersService {
	return &NearbyUsersServiceV2{
		repo:         repo,
		cache:        cache,
		useRedis:     useRedis,
		fallbackToPG: true,
	}
}

func (s *NearbyUsersServiceV2) UpdateUserLocation(
	ctx context.Context,
	userID uuid.UUID,
	deviceID string,
	lat, lon, speed, heading float64,
	isAnonymous bool,
) error {
	location := model.NewUserLocation(userID, deviceID, lat, lon, speed, heading, isAnonymous)

	// 1. Write to PostgreSQL (source of truth)
	if err := s.repo.Upsert(ctx, location); err != nil {
		slog.Error("failed to upsert location in postgres", "error", err)
		return err
	}

	if s.useRedis {
		go s.updateRedisLocation(context.WithoutCancel(ctx), location)
		go s.invalidateNearbyCache(context.WithoutCancel(ctx), userID)
	}

	go func(bgCtx context.Context, uid uuid.UUID, latitude, longitude, spd, hdg float64, devID string) {
		if err := s.repo.SaveHistory(bgCtx, uid, latitude, longitude, spd, hdg, devID); err != nil {
			slog.Error("failed to save location history", slog.Any("error", err), slog.String("user_id", uid.String()))
		}
	}(context.WithoutCancel(ctx), userID, lat, lon, speed, heading, deviceID)

	return nil
}

func (s *NearbyUsersServiceV2) updateRedisLocation(ctx context.Context, loc *model.UserLocation) {
	startTime := time.Now()
	userIDStr := loc.UserID.String()

	if err := s.cache.GeoAdd(ctx, redisGeoKey, loc.Longitude, loc.Latitude, userIDStr); err != nil {
		slog.Warn("failed to add to redis geo", "error", err, "user_id", userIDStr)
		s.redisErrors++
		return
	}

	metaKey := fmt.Sprintf("user_locations:meta:%s", userIDStr)
	metadata := map[string]string{
		"avatar_id":    fmt.Sprintf("%d", loc.AvatarID),
		"color":        loc.Color,
		"speed":        fmt.Sprintf("%f", loc.Speed),
		"heading":      fmt.Sprintf("%f", loc.Heading),
		"is_anonymous": fmt.Sprintf("%t", loc.IsAnonymous),
		"last_update":  loc.LastUpdate.Format(time.RFC3339),
		"device_id":    loc.DeviceID,
	}

	for field, value := range metadata {
		if err := s.cache.HSet(ctx, metaKey, field, value); err != nil {
			slog.Warn("failed to set metadata in redis", "error", err, "user_id", userIDStr)
			s.redisErrors++
			return
		}
	}

	s.cache.HExpire(ctx, metaKey, locationMetadataTTL)

	duration := time.Since(startTime)
	slog.Debug("updated redis location",
		"user_id", userIDStr,
		"duration_ms", duration.Milliseconds())
}

func (s *NearbyUsersServiceV2) invalidateNearbyCache(ctx context.Context, userID uuid.UUID) {
	cacheKey := fmt.Sprintf("nearby_users:%s", userID.String())
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		slog.Debug("failed to invalidate nearby cache", "error", err, "user_id", userID.String())
	}
}

func (s *NearbyUsersServiceV2) GetNearbyUsers(
	ctx context.Context,
	requestingUserID uuid.UUID,
	lat, lon, radiusMeters float64,
) ([]*model.NearbyUser, error) {
	startTime := time.Now()
	userIDStr := requestingUserID.String()

	if s.useRedis {
		if cached := s.getFromCache(ctx, requestingUserID); cached != nil {
			s.cacheHits++
			slog.Debug("[CACHE HIT] nearby users from cache",
				"user_id", userIDStr,
				"count", len(cached),
				"duration_ms", time.Since(startTime).Milliseconds())
			return cached, nil
		}
		s.cacheMisses++
	}

	if s.useRedis {
		users, err := s.getNearbyUsersFromRedis(ctx, requestingUserID, lat, lon, radiusMeters)
		if err == nil {
			go s.cacheNearbyUsers(context.WithoutCancel(ctx), requestingUserID, users)

			duration := time.Since(startTime)
			slog.Debug("[REDIS] nearby users query",
				"user_id", userIDStr,
				"count", len(users),
				"duration_ms", duration.Milliseconds())
			return users, nil
		}

		slog.Warn("[REDIS] failed to get nearby users, falling back to postgres",
			"error", err,
			"user_id", userIDStr)
		s.redisErrors++
	}

	users, err := s.getNearbyUsersFromPostgres(ctx, requestingUserID, lat, lon, radiusMeters)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime)
	slog.Debug("[POSTGRES] nearby users query",
		"user_id", userIDStr,
		"count", len(users),
		"duration_ms", duration.Milliseconds())

	return users, nil
}

func (s *NearbyUsersServiceV2) getFromCache(ctx context.Context, userID uuid.UUID) []*model.NearbyUser {
	cacheKey := fmt.Sprintf("nearby_users:%s", userID.String())
	cached, err := s.cache.Get(ctx, cacheKey)
	if err != nil || cached == "" {
		return nil
	}

	var users []*model.NearbyUser
	if err := json.Unmarshal([]byte(cached), &users); err != nil {
		slog.Debug("failed to unmarshal cached nearby users", "error", err)
		return nil
	}

	return users
}

func (s *NearbyUsersServiceV2) cacheNearbyUsers(ctx context.Context, userID uuid.UUID, users []*model.NearbyUser) {
	cacheKey := fmt.Sprintf("nearby_users:%s", userID.String())
	jsonData, err := json.Marshal(users)
	if err != nil {
		slog.Debug("failed to marshal nearby users for cache", "error", err)
		return
	}

	if err := s.cache.Set(ctx, cacheKey, string(jsonData), nearbyUsersCacheTTLV2); err != nil {
		slog.Debug("failed to cache nearby users", "error", err)
	}
}

func (s *NearbyUsersServiceV2) getNearbyUsersFromRedis(
	ctx context.Context,
	requestingUserID uuid.UUID,
	lat, lon, radiusMeters float64,
) ([]*model.NearbyUser, error) {
	results, err := s.cache.GeoSearchWithDistance(ctx, redisGeoKey, lon, lat, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("redis geo search failed: %w", err)
	}

	if len(results) == 0 {
		return []*model.NearbyUser{}, nil
	}

	nearbyUsers := make([]*model.NearbyUser, 0, len(results))
	usersChan := make(chan *model.NearbyUser, len(results))
	errChan := make(chan error, len(results))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentGoroutines)

	for _, result := range results {
		if result.Member == requestingUserID.String() {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(member string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			user, err := s.fetchUserMetadata(ctx, member, lat, lon)
			if err != nil {
				errChan <- err
				return
			}

			usersChan <- user
		}(result.Member)
	}

	wg.Wait()
	close(usersChan)
	close(errChan)

	for user := range usersChan {
		nearbyUsers = append(nearbyUsers, user)
		if len(nearbyUsers) >= maxNearbyUsersLimit {
			break
		}
	}

	errorCount := len(errChan)
	if errorCount > 0 {
		slog.Debug("some metadata fetches failed",
			"error_count", errorCount,
			"total_results", len(results))
	}

	return nearbyUsers, nil
}

func (s *NearbyUsersServiceV2) fetchUserMetadata(
	ctx context.Context,
	userIDStr string,
	centerLat, centerLon float64,
) (*model.NearbyUser, error) {
	metaKey := fmt.Sprintf("user_locations:meta:%s", userIDStr)
	meta, err := s.cache.HGetAll(ctx, metaKey)
	if err != nil || len(meta) == 0 {
		return nil, fmt.Errorf("metadata not found: %w", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	avatarID, _ := strconv.Atoi(meta["avatar_id"])
	speed, _ := strconv.ParseFloat(meta["speed"], 64)
	heading, _ := strconv.ParseFloat(meta["heading"], 64)
	isAnonymous := meta["is_anonymous"] == "true"
	lastUpdate, _ := time.Parse(time.RFC3339, meta["last_update"])

	privacyLat, privacyLon := model.ApplyPrivacyOffset(centerLat, centerLon)

	return &model.NearbyUser{
		UserID:      userID,
		AnonymousID: model.GenerateAnonymousID(userID),
		Latitude:    privacyLat,
		Longitude:   privacyLon,
		AvatarID:    fmt.Sprintf("avatar_%d", avatarID),
		Color:       meta["color"],
		Speed:       speed,
		Heading:     heading,
		IsAnonymous: isAnonymous,
		LastUpdate:  lastUpdate,
	}, nil
}

func (s *NearbyUsersServiceV2) getNearbyUsersFromPostgres(
	ctx context.Context,
	requestingUserID uuid.UUID,
	lat, lon, radiusMeters float64,
) ([]*model.NearbyUser, error) {
	const maxUsers = 100

	locations, err := s.repo.FindNearbyUsers(ctx, lat, lon, radiusMeters, maxUsers+1)
	if err != nil {
		return nil, err
	}

	nearbyUsers := make([]*model.NearbyUser, 0, len(locations))
	for _, loc := range locations {
		if loc.UserID == requestingUserID {
			continue
		}

		privacyLat, privacyLon := model.ApplyPrivacyOffset(loc.Latitude, loc.Longitude)

		nearbyUser := &model.NearbyUser{
			UserID:      loc.UserID,
			AnonymousID: model.GenerateAnonymousID(loc.UserID),
			Latitude:    privacyLat,
			Longitude:   privacyLon,
			AvatarID:    fmt.Sprintf("avatar_%d", loc.AvatarID),
			Color:       loc.Color,
			Speed:       loc.Speed,
			Heading:     loc.Heading,
			LastUpdate:  loc.LastUpdate,
			IsAnonymous: loc.IsAnonymous,
		}

		nearbyUsers = append(nearbyUsers, nearbyUser)

		if len(nearbyUsers) >= maxUsers {
			break
		}
	}

	return nearbyUsers, nil
}

func (s *NearbyUsersServiceV2) CleanupStaleLocations(ctx context.Context) error {
	return s.repo.DeleteStale(ctx, staleLocationThresholdSeconds)
}

// GetMetrics returns service metrics
func (s *NearbyUsersServiceV2) GetMetrics() map[string]interface{} {
	total := s.cacheHits + s.cacheMisses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(s.cacheHits) / float64(total) * percentageMultiplier
	}

	return map[string]interface{}{
		"cache_hits":     s.cacheHits,
		"cache_misses":   s.cacheMisses,
		"cache_hit_rate": hitRate,
		"redis_errors":   s.redisErrors,
		"use_redis":      s.useRedis,
	}
}
