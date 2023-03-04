package locationtype

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type LocationTypeUseCase interface {
	CreateLocationType(dto CreateLocationTypeDTO) (*LocationTypeDTO, error)
}

type LocationTypeUseCaseImpl struct {
	LocationTypeRepository repository.LocationTypeRepository
}

func NewLocationTypeUseCase(locationTypeRepository repository.LocationTypeRepository) LocationTypeUseCase {
	return &LocationTypeUseCaseImpl{
		LocationTypeRepository: locationTypeRepository,
	}
}

func (l *LocationTypeUseCaseImpl) CreateLocationType(dto CreateLocationTypeDTO) (*LocationTypeDTO, error) {

	locationType, err := entities.NewLocattionType(dto.Name)
	if err != nil {
		return nil, err
	}

	if err := l.LocationTypeRepository.Save(locationType); err != nil {
		return nil, err
	}

	return &LocationTypeDTO{
		ID:   locationType.ID,
		Name: locationType.Name,
	}, nil
}
