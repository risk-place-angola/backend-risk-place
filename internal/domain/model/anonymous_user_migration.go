package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AnonymousUserMigration struct {
	ID                       uuid.UUID
	AnonymousSessionID       uuid.UUID
	DeviceID                 string
	UserID                   uuid.UUID
	AlertsMigrated           int
	SubscriptionsMigrated    int
	SettingsMigrated         bool
	LocationSharingsMigrated int
	MigrationType            string
	StartedAt                time.Time
	CompletedAt              *time.Time
	FailedAt                 *time.Time
	ErrorMessage             *string
}

func NewAnonymousUserMigration(
	anonymousSessionID uuid.UUID,
	deviceID string,
	userID uuid.UUID,
	migrationType string,
) (*AnonymousUserMigration, error) {
	migration := &AnonymousUserMigration{
		ID:                 uuid.New(),
		AnonymousSessionID: anonymousSessionID,
		DeviceID:           deviceID,
		UserID:             userID,
		MigrationType:      migrationType,
		StartedAt:          time.Now(),
	}

	if err := migration.Validate(); err != nil {
		return nil, err
	}

	return migration, nil
}

func (m *AnonymousUserMigration) Validate() error {
	if m.AnonymousSessionID == uuid.Nil {
		return errors.New("anonymous_session_id is required")
	}

	if m.DeviceID == "" {
		return errors.New("device_id is required")
	}

	if m.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	validTypes := map[string]bool{
		"signup": true,
		"login":  true,
		"manual": true,
	}
	if !validTypes[m.MigrationType] {
		return errors.New("migration_type must be 'signup', 'login', or 'manual'")
	}

	return nil
}

func (m *AnonymousUserMigration) MarkCompleted() {
	now := time.Now()
	m.CompletedAt = &now
}

func (m *AnonymousUserMigration) MarkFailed(errorMsg string) {
	now := time.Now()
	m.FailedAt = &now
	m.ErrorMessage = &errorMsg
}

func (m *AnonymousUserMigration) IsCompleted() bool {
	return m.CompletedAt != nil
}

func (m *AnonymousUserMigration) IsFailed() bool {
	return m.FailedAt != nil
}

type DeviceUserMapping struct {
	ID                 uuid.UUID
	DeviceID           string
	AnonymousSessionID uuid.UUID
	UserID             uuid.UUID
	MappedAt           time.Time
	UnmappedAt         *time.Time
	IsActive           bool
}

func NewDeviceUserMapping(
	deviceID string,
	anonymousSessionID uuid.UUID,
	userID uuid.UUID,
) (*DeviceUserMapping, error) {
	mapping := &DeviceUserMapping{
		ID:                 uuid.New(),
		DeviceID:           deviceID,
		AnonymousSessionID: anonymousSessionID,
		UserID:             userID,
		MappedAt:           time.Now(),
		IsActive:           true,
	}

	if err := mapping.Validate(); err != nil {
		return nil, err
	}

	return mapping, nil
}

func (m *DeviceUserMapping) Validate() error {
	if m.DeviceID == "" {
		return errors.New("device_id is required")
	}

	if m.AnonymousSessionID == uuid.Nil {
		return errors.New("anonymous_session_id is required")
	}

	if m.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	return nil
}

func (m *DeviceUserMapping) Deactivate() {
	now := time.Now()
	m.UnmappedAt = &now
	m.IsActive = false
}
