package dto

type DangerZoneDTO struct {
	ID           string  `json:"id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	GridCellID   string  `json:"grid_cell_id"`
	IncidentCount int     `json:"incident_count"`
	RiskScore    float64 `json:"risk_score"`
	RiskLevel    string  `json:"risk_level"`
	CalculatedAt string  `json:"calculated_at"`
}

type GetDangerZonesRequest struct {
	Latitude     float64 `json:"latitude" validate:"required,latitude"`
	Longitude    float64 `json:"longitude" validate:"required,longitude"`
	RadiusMeters float64 `json:"radius_meters" validate:"required,min=100,max=10000"`
}

type GetDangerZonesResponse struct {
	Zones      []DangerZoneDTO `json:"zones"`
	TotalCount int             `json:"total_count"`
}
