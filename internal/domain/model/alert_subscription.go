package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AlertSubscription struct {
	ID                 uuid.UUID
	AlertID            uuid.UUID
	UserID             *uuid.UUID // Nullable - either UserID OR (AnonymousSessionID + DeviceID)
	AnonymousSessionID *uuid.UUID // Set for anonymous users
	DeviceID           *string    // Set for anonymous users
	SubscribedAt       time.Time
}

// IsAnonymous returns true if this subscription was made by an anonymous user
func (s *AlertSubscription) IsAnonymous() bool {
	return s.AnonymousSessionID != nil && s.DeviceID != nil
}

// IsAuthenticated returns true if this subscription was made by an authenticated user
func (s *AlertSubscription) IsAuthenticated() bool {
	return s.UserID != nil
}

// NewAlertSubscription creates a new subscription for an authenticated user
func NewAlertSubscription(alertID, userID uuid.UUID) (*AlertSubscription, error) {
	if alertID == uuid.Nil {
		return nil, errors.New("alert ID is required")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	return &AlertSubscription{
		ID:           uuid.New(),
		AlertID:      alertID,
		UserID:       &userID,
		SubscribedAt: time.Now(),
	}, nil
}

// NewAnonymousAlertSubscription creates a new subscription for an anonymous user
func NewAnonymousAlertSubscription(alertID, anonymousSessionID uuid.UUID, deviceID string) (*AlertSubscription, error) {
	if alertID == uuid.Nil {
		return nil, errors.New("alert ID is required")
	}

	if anonymousSessionID == uuid.Nil {
		return nil, errors.New("anonymous session ID is required")
	}

	if deviceID == "" {
		return nil, errors.New("device ID is required")
	}

	return &AlertSubscription{
		ID:                 uuid.New(),
		AlertID:            alertID,
		AnonymousSessionID: &anonymousSessionID,
		DeviceID:           &deviceID,
		SubscribedAt:       time.Now(),
	}, nil
}
