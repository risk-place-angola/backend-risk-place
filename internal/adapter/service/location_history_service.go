package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// LocationHistoryService manages location history with Redis TTL
type LocationHistoryService struct {
	cache        CacheService
	retentionTTL time.Duration
	enableRedis  bool
}

type LocationHistoryEntry struct {
	UserID    string    `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Timestamp time.Time `json:"timestamp"`
}

type CacheService interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	ZAdd(ctx context.Context, key string, score float64, member string) error
	ZRangeByScore(ctx context.Context, key string, minScore, maxScore float64) ([]string, error)
	ZRemRangeByScore(ctx context.Context, key string, minScore, maxScore float64) error
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

const (
	// Default retention: 7 days (configurable via env)
	defaultRetentionTTL      = 7 * 24 * time.Hour
	locationHistoryKeyPrefix = "location_history:"
	hoursPerDay              = 24
)

func NewLocationHistoryService(cache CacheService, enableRedis bool, retentionDays int) *LocationHistoryService {
	ttl := defaultRetentionTTL
	if retentionDays > 0 {
		ttl = time.Duration(retentionDays) * hoursPerDay * time.Hour
	}

	return &LocationHistoryService{
		cache:        cache,
		retentionTTL: ttl,
		enableRedis:  enableRedis,
	}
}

// SaveHistory saves location to Redis with automatic TTL expiration
// Uses Redis Sorted Set (ZADD) with timestamp as score for efficient time-based queries
func (s *LocationHistoryService) SaveHistory(
	ctx context.Context,
	userID uuid.UUID,
	lat, lon, speed, heading float64,
	deviceID string,
) error {
	if !s.enableRedis {
		// Skip if Redis is disabled
		return nil
	}

	entry := LocationHistoryEntry{
		UserID:    userID.String(),
		DeviceID:  deviceID,
		Latitude:  lat,
		Longitude: lon,
		Speed:     speed,
		Heading:   heading,
		Timestamp: time.Now(),
	}

	// Serialize to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal location entry: %w", err)
	}

	// Redis key: location_history:{user_id}
	key := locationHistoryKeyPrefix + userID.String()
	score := float64(entry.Timestamp.Unix()) // Unix timestamp as score

	// Add to sorted set with timestamp as score
	if err := s.cache.ZAdd(ctx, key, score, string(data)); err != nil {
		return fmt.Errorf("failed to add location to sorted set: %w", err)
	}

	// Set TTL on the key (Redis will auto-delete after retention period)
	if err := s.cache.Expire(ctx, key, s.retentionTTL); err != nil {
		slog.Warn("failed to set TTL on location history", "error", err, "user_id", userID)
	}

	// Auto-cleanup old entries (beyond retention)
	// This is a probabilistic approach: 1% of writes trigger cleanup
	if shouldCleanup() {
		go s.cleanupOldEntries(context.WithoutCancel(ctx), key)
	}

	return nil
}

// GetHistory retrieves location history for a user within time range
func (s *LocationHistoryService) GetHistory(
	ctx context.Context,
	userID uuid.UUID,
	startTime, endTime time.Time,
) ([]LocationHistoryEntry, error) {
	if !s.enableRedis {
		return nil, fmt.Errorf("redis is disabled")
	}

	key := locationHistoryKeyPrefix + userID.String()
	minScore := float64(startTime.Unix())
	maxScore := float64(endTime.Unix())

	// Get entries from sorted set by score range
	results, err := s.cache.ZRangeByScore(ctx, key, minScore, maxScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get location history: %w", err)
	}

	entries := make([]LocationHistoryEntry, 0, len(results))
	for _, data := range results {
		var entry LocationHistoryEntry
		if err := json.Unmarshal([]byte(data), &entry); err != nil {
			slog.Warn("failed to unmarshal location entry", "error", err)
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetRecentHistory retrieves last N hours of location history
func (s *LocationHistoryService) GetRecentHistory(
	ctx context.Context,
	userID uuid.UUID,
	hours int,
) ([]LocationHistoryEntry, error) {
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)
	return s.GetHistory(ctx, userID, startTime, endTime)
}

// cleanupOldEntries removes entries older than retention period
func (s *LocationHistoryService) cleanupOldEntries(ctx context.Context, key string) {
	cutoffTime := time.Now().Add(-s.retentionTTL)
	minScore := float64(0) // Beginning of time
	maxScore := float64(cutoffTime.Unix())

	if err := s.cache.ZRemRangeByScore(ctx, key, minScore, maxScore); err != nil {
		slog.Error("failed to cleanup old location history",
			"error", err,
			"key", key,
			"cutoff", cutoffTime,
		)
	}
}

// shouldCleanup returns true ~1% of the time (probabilistic cleanup)
func shouldCleanup() bool {
	// Using time-based pseudo-random: cleanup every ~100th call
	return time.Now().UnixNano()%100 == 0
}
