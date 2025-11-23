package dto

import (
	"time"

	"github.com/google/uuid"
)

// RiskTypeResponse represents the response structure for a risk type
type RiskTypeResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IconURL       *string   `json:"icon_url,omitempty"`
	DefaultRadius int       `json:"default_radius"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RiskTopicResponse represents the response structure for a risk topic
type RiskTopicResponse struct {
	ID          uuid.UUID `json:"id"`
	RiskTypeID  uuid.UUID `json:"risk_type_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IconURL     *string   `json:"icon_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RiskTypesListResponse represents the response for listing risk types
type RiskTypesListResponse struct {
	Data []RiskTypeResponse `json:"data"`
}

// RiskTopicsListResponse represents the response for listing risk topics
type RiskTopicsListResponse struct {
	Data []RiskTopicResponse `json:"data"`
}
