package placetype

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type PlaceTypeUseCase interface {
	CreatePlaceType(dto CreatePlaceTypeDTO) (*PlaceTypeDTO, error)
	UpdatePlaceType(id string, dto UpdatePlaceTypeDTO) (*PlaceTypeDTO, error)
	FindAllPlaceTypes() ([]*PlaceTypeDTO, error)
	FindByIdPlaceType(id string) (*PlaceTypeDTO, error)
	DeletePlaceType(id string) error
}

type PlaceTypeUseCaseImpl struct {
	PlaceTypeRepository repository.PlaceTypeRepository
}

func NewPlaceTypeUseCase(placeTypeRepository repository.PlaceTypeRepository) PlaceTypeUseCase {
	return &PlaceTypeUseCaseImpl{
		PlaceTypeRepository: placeTypeRepository,
	}
}

func (l *PlaceTypeUseCaseImpl) CreatePlaceType(dto CreatePlaceTypeDTO) (*PlaceTypeDTO, error) {

	placeType, err := entities.NewLocattionType(dto.Name)
	if err != nil {
		return nil, err
	}

	if err := l.PlaceTypeRepository.Save(placeType); err != nil {
		return nil, err
	}

	return &PlaceTypeDTO{
		ID:   placeType.ID,
		Name: placeType.Name,
	}, nil
}

func (l *PlaceTypeUseCaseImpl) UpdatePlaceType(id string, dto UpdatePlaceTypeDTO) (*PlaceTypeDTO, error) {

	placeType, err := l.PlaceTypeRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := placeType.Update(dto.Name); err != nil {
		return nil, err
	}

	if err := l.PlaceTypeRepository.Update(placeType); err != nil {
		return nil, err
	}

	return &PlaceTypeDTO{
		ID:   placeType.ID,
		Name: placeType.Name,
	}, nil
}

func (l *PlaceTypeUseCaseImpl) FindAllPlaceTypes() ([]*PlaceTypeDTO, error) {

	placeTypes, err := l.PlaceTypeRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var placeTypesDTO []*PlaceTypeDTO
	for _, placeType := range placeTypes {
		placeTypesDTO = append(placeTypesDTO, &PlaceTypeDTO{
			ID:   placeType.ID,
			Name: placeType.Name,
		})
	}

	return placeTypesDTO, nil
}

func (l *PlaceTypeUseCaseImpl) FindByIdPlaceType(id string) (*PlaceTypeDTO, error) {

	placeType, err := l.PlaceTypeRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &PlaceTypeDTO{
		ID:   placeType.ID,
		Name: placeType.Name,
	}, nil
}

func (l *PlaceTypeUseCaseImpl) DeletePlaceType(id string) error {

	placeType, err := l.PlaceTypeRepository.FindByID(id)
	if err != nil {
		return err
	}

	return l.PlaceTypeRepository.Delete(placeType.ID)
}
