package model

import (
	"time"

	"github.com/google/uuid"
)

type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusVerified ReportStatus = "verified"
	ReportStatusResolved ReportStatus = "resolved"
	ReportStatusRejected ReportStatus = "rejected"
)

type Report struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RiskTypeID   uuid.UUID
	RiskTopicID  uuid.UUID
	Description  string
	Latitude     float64
	Longitude    float64
	Province     string
	Municipality string
	Neighborhood string
	Address      string
	ImageURL     string
	Status       ReportStatus
	ReviewedBy   uuid.UUID
	ResolvedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
