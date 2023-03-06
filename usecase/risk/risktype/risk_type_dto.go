package risktype

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type RiskTypeDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRiskTypeDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRiskTypeDTO struct {
	CreateRiskTypeDTO
}

func (r *RiskTypeDTO) ToRiskType() *entities.RiskType {
	return &entities.RiskType{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

func (r *RiskTypeDTO) FromRiskType(riskType *entities.RiskType) *RiskTypeDTO {
	return &RiskTypeDTO{
		ID:          riskType.ID,
		Name:        riskType.Name,
		Description: riskType.Description,
	}
}

func (r *RiskTypeDTO) FromRiskTypes(riskTypes []*entities.RiskType) []*RiskTypeDTO {
	var riskTypeDTOs []*RiskTypeDTO
	for _, riskType := range riskTypes {
		riskTypeDTOs = append(riskTypeDTOs, r.FromRiskType(riskType))
	}
	return riskTypeDTOs
}
