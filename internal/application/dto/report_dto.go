package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type ReportDTO struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RiskTypeID   uuid.UUID `json:"risk_type_id"`
	RiskTopicID  uuid.UUID `json:"risk_topic_id,omitempty"`
	Description  string    `json:"description,omitempty"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Province     string    `json:"province,omitempty"`
	Municipality string    `json:"municipality,omitempty"`
	Neighborhood string    `json:"neighborhood,omitempty"`
	Address      string    `json:"address,omitempty"`
	ImageURL     string    `json:"image_url,omitempty"`
	Status       string    `json:"status"`
	ReviewedBy   uuid.UUID `json:"reviewed_by,omitempty"`
	ResolvedAt   time.Time `json:"resolved_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReportCreate struct {
	UserID       string  `json:"user_id" validate:"required,uuid"`
	RiskTypeID   string  `json:"risk_type_id" validate:"required,uuid"`
	RiskTopicID  string  `json:"risk_topic_id"`
	Description  string  `json:"description"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	Province     string  `json:"province,omitempty"`
	Municipality string  `json:"municipality,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	Address      string  `json:"address,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
}

type ReportResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	RiskTypeID   string  `json:"risk_type_id"`
	RiskTopicID  string  `json:"risk_topic_id,omitempty"`
	Description  string  `json:"description"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Province     string  `json:"province,omitempty"`
	Municipality string  `json:"municipality,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	Address      string  `json:"address,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	Status       string  `json:"status"`
	ResolvedAt   string  `json:"resolved_at,omitempty"`
	CreatedAt    string  `json:"created_at"`
	ReviewedBy   string  `json:"reviewed_by,omitempty"`
}

type VerifyReportRequest struct {
	ModeratorID string `json:"moderator_id" validate:"required,uuid"`
}

type ResolveReportRequest struct {
	ModeratorID string `json:"moderator_id" validate:"required,uuid"`
}

func ReportToDTO(r *model.Report) ReportDTO {
	status := string(r.Status)
	return ReportDTO{
		ID:           r.ID,
		UserID:       r.UserID,
		RiskTypeID:   r.RiskTypeID,
		RiskTopicID:  r.RiskTopicID,
		Description:  r.Description,
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
		Province:     r.Province,
		Municipality: r.Municipality,
		Neighborhood: r.Neighborhood,
		Address:      r.Address,
		ImageURL:     r.ImageURL,
		Status:       status,
		ReviewedBy:   r.ReviewedBy,
		ResolvedAt:   r.ResolvedAt,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
