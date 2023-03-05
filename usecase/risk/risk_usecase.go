package risk_usecase

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type RiskUseCase interface {
	CreateRisk(dto CreateRiskDTO) (*RiskDTO, error)
}

type RiskUseCaseImpl struct {
	RiskRepository repository.RiskRepository
}

func NewRiskUseCase(riskRepository repository.RiskRepository) RiskUseCase {
	return &RiskUseCaseImpl{
		RiskRepository: riskRepository,
	}
}

func (r *RiskUseCaseImpl) CreateRisk(dto CreateRiskDTO) (*RiskDTO, error) {

	riskEntity := &entities.Risk{
		RiskTypeID:     dto.RiskTypeID,
		LocationTypeID: dto.LocationTypeID,
		Name:           dto.Name,
		Latitude:       dto.Latitude,
		Longitude:      dto.Longitude,
		Description:    dto.Description,
	}

	risk, err := entities.NewRisk(riskEntity)
	if err != nil {
		return nil, err
	}

	if err := r.RiskRepository.Save(risk); err != nil {
		return nil, err
	}

	dtoRisk := &RiskDTO{}

	return dtoRisk.FromRisk(risk), nil

}
