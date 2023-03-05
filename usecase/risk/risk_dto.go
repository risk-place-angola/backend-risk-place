package risk_usecase

type CreateRiskDTO struct {
	RiskTypeID     string  `json:"risk_type_id"`
	LocationTypeID string  `json:"location_type_id"`
	Name           string  `json:"name"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Description    string  `json:"description"`
}

type UpdateRiskDTO struct {
	*CreateRiskDTO
}

type RiskDTO struct {
	ID             string  `json:"id"`
	RiskTypeID     string  `json:"risk_type_id"`
	LocationTypeID string  `json:"location_type_id"`
	Name           string  `json:"name"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Description    string  `json:"description"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
