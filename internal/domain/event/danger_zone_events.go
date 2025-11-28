package event

import (
	"github.com/google/uuid"
)

type DangerZoneEnteredEvent struct {
	UserID       uuid.UUID
	DeviceID     string
	Latitude     float64
	Longitude    float64
	ZoneID       uuid.UUID
	RiskLevel    string
	IncidentCount int
}

func (e DangerZoneEnteredEvent) Name() string {
	return "danger_zone.entered"
}
