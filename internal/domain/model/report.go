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
	ID                uuid.UUID
	UserID            uuid.UUID
	RiskTypeID        uuid.UUID
	RiskTypeName      string
	RiskTypeIconPath  *string
	RiskTopicID       uuid.UUID
	RiskTopicName     string
	RiskTopicIconPath *string
	Description       string
	Latitude          float64
	Longitude         float64
	Province          string
	Municipality      string
	Neighborhood      string
	Address           string
	ImageURL          string
	Status            ReportStatus
	ReviewedBy        uuid.UUID
	ResolvedAt        time.Time
	VerificationCount int
	RejectionCount    int
	ExpiresAt         *time.Time
	IsPrivate         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
