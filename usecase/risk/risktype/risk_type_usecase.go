package risktype

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type RiskTypeUseCase interface {
	CreateRiskType(dto *CreateRiskTypeDTO) (*RiskTypeDTO, error)
	UpdateRiskType(id string, dto *UpdateRiskTypeDTO) (*RiskTypeDTO, error)
	FindAllRiskTypes() ([]*RiskTypeDTO, error)
	FindRiskTypeByID(id string) (*RiskTypeDTO, error)
	RemoveRiskType(id string) error
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

func (r *RiskTypeUseCaseImpl) UpdateRiskType(id string, dto *UpdateRiskTypeDTO) (*RiskTypeDTO, error) {

	risktype, err := r.RiskTypeRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := risktype.Update(dto.Name, dto.Description); err != nil {
		return nil, err
	}

	if err := r.RiskTypeRepository.Update(risktype); err != nil {
		return nil, err
	}

	dtoRiskType := &RiskTypeDTO{}

	return dtoRiskType.FromRiskType(risktype), nil
}

func (r *RiskTypeUseCaseImpl) FindAllRiskTypes() ([]*RiskTypeDTO, error) {

	riskTypes, err := r.RiskTypeRepository.FindAll()
	if err != nil {
		return nil, err
	}

	dtoRiskType := &RiskTypeDTO{}

	return dtoRiskType.FromRiskTypes(riskTypes), nil

}

func (r *RiskTypeUseCaseImpl) FindRiskTypeByID(id string) (*RiskTypeDTO, error) {
	risktype, err := r.RiskTypeRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	dtoRiskType := &RiskTypeDTO{}

	return dtoRiskType.FromRiskType(risktype), nil
}

func (r *RiskTypeUseCaseImpl) RemoveRiskType(id string) error {
	risktype, err := r.RiskTypeRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err := r.RiskTypeRepository.Delete(risktype.ID); err != nil {
		return err
	}

	return nil
}
