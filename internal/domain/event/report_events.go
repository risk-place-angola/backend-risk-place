package event

import "github.com/google/uuid"

type ReportVerifiedEvent struct {
	ReportID uuid.UUID
	UserID   uuid.UUID
}

func (r ReportVerifiedEvent) Name() string {
	return "ReportVerified"
}

type ReportCreatedEvent struct {
	ReportID  uuid.UUID
	UserID    []uuid.UUID
	Message   string
	Latitude  float64
	Longitude float64
	Radius    float64
	RiskType  string
}

func (e ReportCreatedEvent) Name() string { return "ReportCreated" }

type ReportResolvedEvent struct {
	ReportID uuid.UUID
	Message  string
	UserIDs  []string
}

func (e ReportResolvedEvent) Name() string { return "ReportResolved" }
