package risktype

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type RiskTypeUseCase interface {
	CreateRiskType(dto *CreateRiskTypeDTO) (*RiskTypeDTO, error)
}

type RiskTypeUseCaseImpl struct {
	RiskTypeRepository repository.RiskTypeRepository
}

func NewRiskTypeUseCase(riskTypeRepository repository.RiskTypeRepository) RiskTypeUseCase {
	return &RiskTypeUseCaseImpl{
		RiskTypeRepository: riskTypeRepository,
	}
}

func (r *RiskTypeUseCaseImpl) CreateRiskType(dto *CreateRiskTypeDTO) (*RiskTypeDTO, error) {

	riskType, err := entities.NewRiskType(dto.Name, dto.Description)
	if err != nil {
		return nil, err
	}

	if err := r.RiskTypeRepository.Save(riskType); err != nil {
		return nil, err
	}

	dtoRiskType := &RiskTypeDTO{}

	return dtoRiskType.FromRiskType(riskType), nil
}
