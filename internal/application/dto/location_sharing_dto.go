package dto

import (
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type CreateLocationSharingRequest struct {
	DurationMinutes int     `json:"duration_minutes" validate:"required,min=1,max=1440"`
	Latitude        float64 `json:"latitude"         validate:"required,latitude"`
	Longitude       float64 `json:"longitude"        validate:"required,longitude"`
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude"  validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type LocationSharingResponse struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	ShareLink string    `json:"share_link"`
	ExpiresAt time.Time `json:"expires_at"`
}

type PublicLocationResponse struct {
	UserName    string    `json:"user_name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	LastUpdated time.Time `json:"last_updated"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
}

func ToLocationSharingResponse(sharing *model.LocationSharing, baseURL string) *LocationSharingResponse {
	return &LocationSharingResponse{
		ID:        sharing.ID.String(),
		Token:     sharing.Token,
		ShareLink: baseURL + "/share/" + sharing.Token,
		ExpiresAt: sharing.ExpiresAt,
	}
}

func ToPublicLocationResponse(sharing *model.LocationSharing) *PublicLocationResponse {
	return &PublicLocationResponse{
		UserName:    sharing.OwnerName,
		Latitude:    sharing.Latitude,
		Longitude:   sharing.Longitude,
		LastUpdated: sharing.LastUpdatedAt,
		ExpiresAt:   sharing.ExpiresAt,
		IsActive:    sharing.IsActive && !sharing.IsExpired(),
	}
}
