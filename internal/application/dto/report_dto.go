package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type ReportDTO struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	RiskTypeID        uuid.UUID  `json:"risk_type_id"`
	RiskTypeName      string     `json:"risk_type_name,omitempty"`
	RiskTypeIconURL   *string    `json:"risk_type_icon_url,omitempty"`
	RiskTopicID       uuid.UUID  `json:"risk_topic_id,omitempty"`
	RiskTopicName     string     `json:"risk_topic_name,omitempty"`
	RiskTopicIconURL  *string    `json:"risk_topic_icon_url,omitempty"`
	Description       string     `json:"description,omitempty"`
	Latitude          float64    `json:"latitude"`
	Longitude         float64    `json:"longitude"`
	Province          string     `json:"province,omitempty"`
	Municipality      string     `json:"municipality,omitempty"`
	Neighborhood      string     `json:"neighborhood,omitempty"`
	Address           string     `json:"address,omitempty"`
	ImageURL          string     `json:"image_url,omitempty"`
	Status            string     `json:"status"`
	ReviewedBy        uuid.UUID  `json:"reviewed_by,omitempty"`
	ResolvedAt        time.Time  `json:"resolved_at,omitempty"`
	VerificationCount int        `json:"verification_count"`
	RejectionCount    int        `json:"rejection_count"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ReportCreate struct {
	UserID       string  `json:"user_id"`
	RiskTypeID   string  `json:"risk_type_id"`
	RiskTopicID  string  `json:"risk_topic_id"`
	Description  string  `json:"description"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
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

type UpdateReportLocationRequest struct {
	Latitude     float64 `json:"latitude"               validate:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude"              validate:"required,min=-180,max=180"`
	Address      string  `json:"address,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	Municipality string  `json:"municipality,omitempty"`
	Province     string  `json:"province,omitempty"`
}

type UpdateReportLocationResponse struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	UpdatedAt string `json:"updated_at"`
}

// ListReportsQueryParams represents query parameters for listing reports with pagination
type ListReportsQueryParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status,omitempty"`
	Sort   string `json:"sort,omitempty"`
	Order  string `json:"order,omitempty"`
}

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasMore     bool `json:"has_more"`
	HasPrevious bool `json:"has_previous"`
}

// ListReportsResponse represents the response for listing reports with pagination
type ListReportsResponse struct {
	Reports    []ReportDTO        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// NearbyReportsQueryParams represents query parameters for nearby reports
type NearbyReportsQueryParams struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	Limit     int     `json:"limit,omitempty"`
}

// ReportWithDistance extends ReportDTO with distance information
type ReportWithDistance struct {
	ReportDTO
	Distance float64 `json:"distance"` // Distance in meters
}

// NearbyReportsResponse represents the response for nearby reports
type NearbyReportsResponse struct {
	Reports []ReportWithDistance `json:"data"`
}

func ReportToDTO(r *model.Report) ReportDTO {
	status := string(r.Status)

	var riskTypeIconURL *string
	if r.RiskTypeIconPath != nil && *r.RiskTypeIconPath != "" {
		url := "/api/v1/storage/" + *r.RiskTypeIconPath
		riskTypeIconURL = &url
	}

	var riskTopicIconURL *string
	if r.RiskTopicIconPath != nil && *r.RiskTopicIconPath != "" {
		url := "/api/v1/storage/" + *r.RiskTopicIconPath
		riskTopicIconURL = &url
	}

	return ReportDTO{
		ID:                r.ID,
		UserID:            r.UserID,
		RiskTypeID:        r.RiskTypeID,
		RiskTypeName:      r.RiskTypeName,
		RiskTypeIconURL:   riskTypeIconURL,
		RiskTopicID:       r.RiskTopicID,
		RiskTopicName:     r.RiskTopicName,
		RiskTopicIconURL:  riskTopicIconURL,
		Description:       r.Description,
		Latitude:          r.Latitude,
		Longitude:         r.Longitude,
		Province:          r.Province,
		Municipality:      r.Municipality,
		Neighborhood:      r.Neighborhood,
		Address:           r.Address,
		ImageURL:          r.ImageURL,
		Status:            status,
		ReviewedBy:        r.ReviewedBy,
		ResolvedAt:        r.ResolvedAt,
		VerificationCount: r.VerificationCount,
		RejectionCount:    r.RejectionCount,
		ExpiresAt:         r.ExpiresAt,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

type VoteReportRequest struct {
	VoteType string `json:"vote_type" validate:"required,oneof=upvote downvote"`
}

type VoteReportResponse struct {
	ReportID          string `json:"report_id"`
	VoteType          string `json:"vote_type"`
	VerificationCount int    `json:"verification_count"`
	RejectionCount    int    `json:"rejection_count"`
}

func ReportToDTOWithDistance(r *model.Report, distance float64) ReportWithDistance {
	return ReportWithDistance{
		ReportDTO: ReportToDTO(r),
		Distance:  distance,
	}
}
