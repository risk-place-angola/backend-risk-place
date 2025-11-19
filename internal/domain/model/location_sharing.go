package model

import (
	"time"

	"github.com/google/uuid"
)

type LocationSharing struct {
	ID                 uuid.UUID
	UserID             *uuid.UUID
	AnonymousSessionID *uuid.UUID
	DeviceID           *string
	OwnerName          string
	Token              string
	Latitude           float64
	Longitude          float64
	DurationMinutes    int
	ExpiresAt          time.Time
	LastUpdatedAt      time.Time
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewLocationSharing(latitude, longitude float64, durationMinutes int, ownerName string) *LocationSharing {
	now := time.Now()
	token := uuid.New().String()

	return &LocationSharing{
		ID:              uuid.New(),
		Token:           token,
		OwnerName:       ownerName,
		Latitude:        latitude,
		Longitude:       longitude,
		DurationMinutes: durationMinutes,
		ExpiresAt:       now.Add(time.Duration(durationMinutes) * time.Minute),
		LastUpdatedAt:   now,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (ls *LocationSharing) SetAuthenticatedUser(userID uuid.UUID) {
	ls.UserID = &userID
	ls.AnonymousSessionID = nil
	ls.DeviceID = nil
}

func (ls *LocationSharing) SetAnonymousUser(sessionID uuid.UUID, deviceID string) {
	ls.UserID = nil
	ls.AnonymousSessionID = &sessionID
	ls.DeviceID = &deviceID
}

func (ls *LocationSharing) UpdateLocation(latitude, longitude float64) {
	ls.Latitude = latitude
	ls.Longitude = longitude
	ls.LastUpdatedAt = time.Now()
	ls.UpdatedAt = time.Now()
}

func (ls *LocationSharing) Deactivate() {
	ls.IsActive = false
	ls.UpdatedAt = time.Now()
}

func (ls *LocationSharing) IsExpired() bool {
	return time.Now().After(ls.ExpiresAt)
}

func (ls *LocationSharing) IsValid() bool {
	return ls.IsActive && !ls.IsExpired()
}

func (ls *LocationSharing) IsOwnedByUser(userID uuid.UUID) bool {
	return ls.UserID != nil && *ls.UserID == userID
}

func (ls *LocationSharing) IsOwnedByDevice(deviceID string) bool {
	return ls.DeviceID != nil && *ls.DeviceID == deviceID
}

func (ls *LocationSharing) GetOwnerIdentifier() string {
	if ls.UserID != nil {
		return ls.UserID.String()
	}
	if ls.DeviceID != nil {
		return *ls.DeviceID
	}
	return ""
}
