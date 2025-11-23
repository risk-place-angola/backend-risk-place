package model

import (
	"time"

	"github.com/google/uuid"
)

type EntityType string

const (
	EntityTypeERCE  EntityType = "erce"
	EntityTypeERFCE EntityType = "erfce"
)

type Entity struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	EntityType   EntityType `json:"entity_type"`
	Province     string     `json:"province"`
	Municipality string     `json:"municipality"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	ContactEmail *string    `json:"contact_email,omitempty"`
	ContactPhone *string    `json:"contact_phone,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
