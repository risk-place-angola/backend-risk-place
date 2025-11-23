package dto

import "time"

type SafeRouteRequest struct {
	OriginLat      float64 `json:"origin_lat"      validate:"required,latitude"`
	OriginLon      float64 `json:"origin_lon"      validate:"required,longitude"`
	DestinationLat float64 `json:"destination_lat" validate:"required,latitude"`
	DestinationLon float64 `json:"destination_lon" validate:"required,longitude"`
	MaxRoutes      int     `json:"max_routes"      validate:"omitempty,min=1,max=3"`
}

type WaypointDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Sequence  int     `json:"sequence"`
}

type IncidentDTO struct {
	ReportID     string    `json:"report_id"`
	RiskType     string    `json:"risk_type"`
	RiskTopic    string    `json:"risk_topic"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	DistanceKm   float64   `json:"distance_km"`
	CreatedAt    time.Time `json:"created_at"`
	DaysAgo      int       `json:"days_ago"`
	WeightFactor float64   `json:"weight_factor"`
}

type SafeRouteResponse struct {
	ID                string        `json:"id"`
	OriginLat         float64       `json:"origin_lat"`
	OriginLon         float64       `json:"origin_lon"`
	DestinationLat    float64       `json:"destination_lat"`
	DestinationLon    float64       `json:"destination_lon"`
	Waypoints         []WaypointDTO `json:"waypoints"`
	DistanceKm        float64       `json:"distance_km"`
	EstimatedDuration int           `json:"estimated_duration_minutes"`
	SafetyScore       float64       `json:"safety_score"`
	RiskLevel         string        `json:"risk_level"`
	IncidentCount     int           `json:"incident_count"`
	Incidents         []IncidentDTO `json:"incidents"`
	CalculatedAt      time.Time     `json:"calculated_at"`
}

type HeatmapRequest struct {
	NorthEastLat float64 `json:"north_east_lat" validate:"required,latitude"`
	NorthEastLon float64 `json:"north_east_lon" validate:"required,longitude"`
	SouthWestLat float64 `json:"south_west_lat" validate:"required,latitude"`
	SouthWestLon float64 `json:"south_west_lon" validate:"required,longitude"`
	StartDate    string  `json:"start_date"     validate:"omitempty"`
	EndDate      string  `json:"end_date"       validate:"omitempty"`
	RiskTypeID   string  `json:"risk_type_id"   validate:"omitempty,uuid"`
}

type HeatmapPointDTO struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Weight       float64 `json:"weight"`
	IncidentType string  `json:"incident_type"`
	ReportCount  int     `json:"report_count"`
}

type HeatmapResponse struct {
	Points     []HeatmapPointDTO `json:"points"`
	TotalCount int               `json:"total_count"`
	BoundsInfo struct {
		NorthEastLat float64 `json:"north_east_lat"`
		NorthEastLon float64 `json:"north_east_lon"`
		SouthWestLat float64 `json:"south_west_lat"`
		SouthWestLon float64 `json:"south_west_lon"`
	} `json:"bounds_info"`
}
