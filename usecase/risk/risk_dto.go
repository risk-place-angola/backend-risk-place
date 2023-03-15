package risk_usecase

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type CreateRiskDTO struct {
	RiskTypeID  string  `json:"risk_type_id"`
	PlaceTypeID string  `json:"place_type_id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type UpdateRiskDTO struct {
	CreateRiskDTO
}

type RiskDTO struct {
	ID          string  `json:"id"`
	RiskTypeID  string  `json:"risk_type_id"`
	PlaceTypeID string  `json:"place_type_id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (u *UpdateRiskDTO) ToRiskUpdate() *entities.Risk {
	return &entities.Risk{
		RiskTypeID:  u.RiskTypeID,
		PlaceTypeID: u.PlaceTypeID,
		Name:        u.Name,
		Latitude:    u.Latitude,
		Longitude:   u.Longitude,
		Description: u.Description,
	}
}

func (r *RiskDTO) ToRisk() *entities.Risk {
	return &entities.Risk{
		ID:          r.ID,
		RiskTypeID:  r.RiskTypeID,
		PlaceTypeID: r.PlaceTypeID,
		Name:        r.Name,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Description: r.Description,
	}
}

func (r *RiskDTO) FromRisk(risk *entities.Risk) *RiskDTO {
	r.ID = risk.ID
	r.RiskTypeID = risk.RiskTypeID
	r.PlaceTypeID = risk.PlaceTypeID
	r.Name = risk.Name
	r.Latitude = risk.Latitude
	r.Longitude = risk.Longitude
	r.Description = risk.Description
	r.CreatedAt = risk.CreatedAt.String()
	r.UpdatedAt = risk.UpdatedAt.String()
	return r
}

func (r *RiskDTO) FromRiskList(risks []*entities.Risk) []*RiskDTO {
	var riskDTOs []*RiskDTO
	for _, risk := range risks {
		riskDTOs = append(riskDTOs, r.FromRisk(risk))
	}
	return riskDTOs
}

func NewRiskDTO() *RiskDTO {
	return &RiskDTO{}
}
