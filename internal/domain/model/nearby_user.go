package model

import (
	"math"
	mathrand "math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	hashMultiplier      = 31
	anonymousIDLength   = 5
	totalAvatars        = 21
	privacyOffsetMeters = 75.0
	metersPerDegree     = 111320.0
	twoPi               = 2 * math.Pi
	degreesToRadians    = math.Pi / 180
)

type NearbyUser struct {
	UserID      uuid.UUID `json:"user_id"`
	AnonymousID string    `json:"anonymous_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	AvatarID    string    `json:"avatar_id"`
	Color       string    `json:"color"`
	Speed       float64   `json:"speed"`
	Heading     float64   `json:"heading"`
	LastUpdate  time.Time `json:"last_update"`
	IsAnonymous bool      `json:"is_anonymous"`
}

type UserLocation struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	DeviceID    string    `json:"device_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Speed       float64   `json:"speed"`
	Heading     float64   `json:"heading"`
	AvatarID    int       `json:"avatar_id"`
	Color       string    `json:"color"`
	IsAnonymous bool      `json:"is_anonymous"`
	LastUpdate  time.Time `json:"last_update"`
	CreatedAt   time.Time `json:"created_at"`
}

func getAvatarColors() []string {
	return []string{
		"#4A90E2", "#E94B3C", "#50C878", "#F5A623", "#9B59B6",
		"#E91E63", "#00BCD4", "#FF9800", "#795548", "#607D8B",
		"#FFEB3B", "#8BC34A", "#673AB7", "#009688", "#FF5722",
	}
}

func GenerateAnonymousID(userID uuid.UUID) string {
	hash := 0
	for _, b := range userID.String() {
		hash = hash*hashMultiplier + int(b)
	}
	// Ensure hash is positive for consistent seed
	if hash < 0 {
		hash = -hash
	}
	return "neter_" + generateRandomString(anonymousIDLength, hash)
}

func AssignAvatar(userID uuid.UUID) (int, string) {
	hash := 0
	for _, b := range userID.String() {
		hash = hash*hashMultiplier + int(b)
	}
	
	// Ensure hash is positive for array indexing
	if hash < 0 {
		hash = -hash
	}
	
	avatarID := (hash % totalAvatars) + 1
	colors := getAvatarColors()
	colorIndex := hash % len(colors)
	color := colors[colorIndex]
	return avatarID, color
}

func generateRandomString(length int, seed int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := mathrand.New(mathrand.NewSource(int64(seed))) // #nosec G404 - deterministic randomness needed for consistent IDs
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func NewUserLocation(userID uuid.UUID, deviceID string, lat, lon, speed, heading float64, isAnonymous bool) *UserLocation {
	avatarID, color := AssignAvatar(userID)
	now := time.Now()

	return &UserLocation{
		ID:          uuid.New(),
		UserID:      userID,
		DeviceID:    deviceID,
		Latitude:    lat,
		Longitude:   lon,
		Speed:       speed,
		Heading:     heading,
		AvatarID:    avatarID,
		Color:       color,
		IsAnonymous: isAnonymous,
		LastUpdate:  now,
		CreatedAt:   now,
	}
}

func ApplyPrivacyOffset(lat, lon float64) (float64, float64) {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano())) // #nosec G404 - privacy offset doesn't require cryptographic randomness
	angle := r.Float64() * twoPi
	distance := r.Float64() * privacyOffsetMeters

	latOffset := (distance * math.Cos(angle)) / metersPerDegree
	lonOffset := (distance * math.Sin(angle)) / (metersPerDegree * math.Cos(lat*degreesToRadians))

	return lat + latOffset, lon + lonOffset
}
