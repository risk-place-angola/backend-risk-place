package model

import (
	"time"

	"github.com/google/uuid"
)

type RiskType struct {
	ID                  uuid.UUID
	Name                string
	Description         string
	IconPath            *string
	DefaultRadiusMeters int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type RiskTopic struct {
	ID          uuid.UUID
	RiskTypeID  uuid.UUID
	Name        string
	Description *string
	IconPath    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
