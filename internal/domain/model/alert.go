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
	ID           uuid.UUID
	CreatedBy    uuid.UUID
	RiskTypeID   uuid.UUID
	RiskTopicID  uuid.UUID
	Message      string
	Latitude     float64
	Longitude    float64
	Province     string
	Municipality string
	Neighborhood string
	Address      string
	RadiusMeters int
	Status       AlertStatus
	Severity     Severity
	CreatedAt    time.Time
	ExpiresAt    time.Time
	ResolvedAt   time.Time
}

type Notification struct {
	ID          uuid.UUID
	Type        string
	ReferenceID uuid.UUID
	UserID      uuid.UUID
	SentAt      time.Time
	SeenAt      time.Time
}
