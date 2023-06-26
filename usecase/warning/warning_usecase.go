package warning_usecase

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type IWarningUseCase interface {
	CreateWarning(dto *CreateWarningDTO) (*CreateWarningDTO, error)
	UpdateWarning(id string, dto *UpdateWarningDTO) (*UpdateWarningDTO, error)
	FindAllWarning() ([]*DTO, error)
	FindWarningByID(id string) (*DTO, error)
	RemoveWarning(id string) error
}

type WarningUseCaseImpl struct {
	WarningRepository repository.IWaringRepository
}

func NewWarningUseCase(warningRepository repository.IWaringRepository) IWarningUseCase {
	return &WarningUseCaseImpl{
		WarningRepository: warningRepository,
	}
}

func (w WarningUseCaseImpl) CreateWarning(dto *CreateWarningDTO) (*CreateWarningDTO, error) {
	location := dto.convertLocationStringToFloat()
	entity := &entities.Warning{
		ReportedBy: dto.ReportedBy,
		Fact:       dto.Fact,
		Latitude:   location.Latitude,
		Longitude:  location.Longitude,
		EventState: entities.Pending,
	}
	warning, err := entities.NewWarning(entity)
	if err != nil {
		return nil, err
	}
	err = w.WarningRepository.Save(warning)
	if err != nil {
		return nil, err
	}
	dtoWarning := &DTO{}

	return dtoWarning.FromCreateWarning(warning), nil
}

func (w WarningUseCaseImpl) UpdateWarning(id string, dto *UpdateWarningDTO) (*UpdateWarningDTO, error) {
	warning, err := w.WarningRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	warningDTO := dto.ToWarningUpdate()

	warning.ReportedBy = warningDTO.ReportedBy
	if err := warning.Update(warningDTO); err != nil {
		return nil, err
	}

	err = w.WarningRepository.Update(warning)
	if err != nil {
		return nil, err
	}

	dtoWarning := &DTO{}
	return dtoWarning.FromUpdateWarning(warning), nil
}

func (w WarningUseCaseImpl) FindAllWarning() ([]*DTO, error) {
	warnings, err := w.WarningRepository.FindAll()
	if err != nil {
		return nil, err
	}
	var dtoWarnings *DTO
	return dtoWarnings.FromWarnings(warnings), nil
}

func (w WarningUseCaseImpl) FindWarningByID(id string) (*DTO, error) {
	warning, err := w.WarningRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	dtoWarning := &DTO{}
	return dtoWarning.FromWarning(warning), nil
}

func (w WarningUseCaseImpl) RemoveWarning(id string) error {
	err := w.WarningRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
