package place_usecase

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type PlaceUseCase interface {
	CreatePlace(dto CreatePlaceDTO) (*PlaceDTO, error)
	UpdatePlace(id string, dto UpdatePlaceDTO) (*PlaceDTO, error)
	FindAllPlace() ([]*PlaceDTO, error)
	FindPlaceByID(id string) (*PlaceDTO, error)
}

type PlaceUseCaseImpl struct {
	PlaceRepository repository.PlaceRepository
}

func NewPlaceUseCase(placeRepository repository.PlaceRepository) PlaceUseCase {
	return &PlaceUseCaseImpl{
		PlaceRepository: placeRepository,
	}
}

func (r *PlaceUseCaseImpl) CreatePlace(dto CreatePlaceDTO) (*PlaceDTO, error) {

	placeEntity := &entities.Place{
		RiskTypeID:  dto.RiskTypeID,
		PlaceTypeID: dto.PlaceTypeID,
		Name:        dto.Name,
		Latitude:    dto.Latitude,
		Longitude:   dto.Longitude,
		Description: dto.Description,
	}

	place, err := entities.NewPlace(placeEntity)
	if err != nil {
		return nil, err
	}

	if err := r.PlaceRepository.Save(place); err != nil {
		return nil, err
	}

	dtoPlace := &PlaceDTO{}

	return dtoPlace.FromPlace(place), nil

}

func (r *PlaceUseCaseImpl) UpdatePlace(id string, dto UpdatePlaceDTO) (*PlaceDTO, error) {

	place, err := r.PlaceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	place = dto.ToPlaceUpdate()
	if err := place.Update(place); err != nil {
		return nil, err
	}

	if err := r.PlaceRepository.Update(place); err != nil {
		return nil, err
	}

	dtoPlace := &PlaceDTO{}

	return dtoPlace.FromPlace(place), nil

}

func (r *PlaceUseCaseImpl) FindAllPlace() ([]*PlaceDTO, error) {

	place, err := r.PlaceRepository.FindAll()
	if err != nil {
		return nil, err
	}

	dtoPlace := &PlaceDTO{}
	dtoPlaces := dtoPlace.FromPlaceList(place)

	return dtoPlaces, nil

}

func (r *PlaceUseCaseImpl) FindPlaceByID(id string) (*PlaceDTO, error) {

	place, err := r.PlaceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	dtoPlace := &PlaceDTO{}

	return dtoPlace.FromPlace(place), nil

}
