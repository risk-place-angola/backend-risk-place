package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AnonymousSession struct {
	ID                uuid.UUID
	DeviceID          string
	DeviceFCMToken    string
	DevicePlatform    string
	DeviceModel       string
	Latitude          float64
	Longitude         float64
	AlertRadiusMeters int
	DeviceLanguage    string
	MigratedToUserID  *uuid.UUID
	MigratedAt        *time.Time
	IsActive          bool
	LastSeen          time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

const (
	DefaultAnonymousAlertRadiusMeters = 1000
	MinDeviceIDLength                 = 16
)

func NewAnonymousSession(deviceID, fcmToken, platform, model string) (*AnonymousSession, error) {
	session := &AnonymousSession{
		ID:                uuid.New(),
		DeviceID:          deviceID,
		DeviceFCMToken:    fcmToken,
		DevicePlatform:    platform,
		DeviceModel:       model,
		AlertRadiusMeters: DefaultAnonymousAlertRadiusMeters,
		DeviceLanguage:    "pt",
		IsActive:          true,
		LastSeen:          time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := session.Validate(); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AnonymousSession) Validate() error {
	if s.DeviceID == "" {
		return errors.New("device_id is required")
	}

	if len(s.DeviceID) < MinDeviceIDLength {
		return errors.New("device_id must be at least 16 characters")
	}

	if s.DevicePlatform != "" {
		validPlatforms := map[string]bool{
			"ios":     true,
			"android": true,
			"web":     true,
		}
		if !validPlatforms[s.DevicePlatform] {
			return errors.New("device_platform must be 'ios', 'android', or 'web'")
		}
	}

	if s.AlertRadiusMeters < 0 {
		return errors.New("alert_radius_meters must be positive")
	}

	return nil
}

func (s *AnonymousSession) UpdateLocation(lat, lon float64) {
	s.Latitude = lat
	s.Longitude = lon
	s.LastSeen = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *AnonymousSession) UpdateFCMToken(token string) {
	s.DeviceFCMToken = token
	s.LastSeen = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *AnonymousSession) TouchLastSeen() {
	s.LastSeen = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *AnonymousSession) MigrateTo(userID uuid.UUID) {
	s.MigratedToUserID = &userID
	now := time.Now()
	s.MigratedAt = &now
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

func (s *AnonymousSession) IsMigrated() bool {
	return s.MigratedToUserID != nil
}
