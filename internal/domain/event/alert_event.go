package event

import "github.com/google/uuid"

type AlertCreatedEvent struct {
	AlertID   uuid.UUID
	UserID    []uuid.UUID
	Message   string
	Latitude  float64
	Longitude float64
	Radius    float64
	RiskType  string
	Severity  string
}

func (e AlertCreatedEvent) Name() string { return "AlertCreated" }
