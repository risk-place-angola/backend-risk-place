package model

import (
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusExpired  AlertStatus = "expired"
)

type Alert struct {
	ID                 uuid.UUID
	CreatedBy          *uuid.UUID // Nullable - either CreatedBy OR (AnonymousSessionID + DeviceID)
	AnonymousSessionID *uuid.UUID // Set for anonymous users
	DeviceID           *string    // Set for anonymous users
	RiskTypeID         uuid.UUID
	RiskTopicID        uuid.UUID
	Message            string
	Latitude           float64
	Longitude          float64
	Province           string
	Municipality       string
	Neighborhood       string
	Address            string
	RadiusMeters       int
	Status             AlertStatus
	Severity           Severity
	CreatedAt          time.Time
	ExpiresAt          time.Time
	ResolvedAt         time.Time
}

// IsAnonymous returns true if this alert was created by an anonymous user
func (a *Alert) IsAnonymous() bool {
	return a.AnonymousSessionID != nil && a.DeviceID != nil
}

// IsAuthenticated returns true if this alert was created by an authenticated user
func (a *Alert) IsAuthenticated() bool {
	return a.CreatedBy != nil
}

type Notification struct {
	ID          uuid.UUID
	Type        string
	ReferenceID uuid.UUID
	UserID      uuid.UUID
	SentAt      time.Time
	SeenAt      time.Time
}
