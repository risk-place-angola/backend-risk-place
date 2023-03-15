package place_usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	place_usecase "github.com/risk-place-angola/backend-risk-place/usecase/place"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().Save(gomock.Any()).Return(nil)

	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	place, err := placeUseCase.CreatePlace(place_usecase.CreatePlaceDTO{
		RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
		PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:        "Rangel rua da Lama",
		Latitude:    8.825248,
		Longitude:   13.263879,
		Description: "Risco de inundação",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", place.Name)
}

func TestUpdatePlace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Place{
		ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
		PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:        "Rangel rua da Lama",
		Latitude:    8.825248,
		Longitude:   13.263879,
		Description: "Risco de inundação",
	}

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockPlaceRepository.EXPECT().Update(gomock.Any()).Return(nil)

	updatePlaceDTO := &place_usecase.UpdatePlaceDTO{}
	updatePlaceDTO.RiskTypeID = "99bada49-09d0-4f13-b310-6f8633b38dfe"
	updatePlaceDTO.PlaceTypeID = "dd3aadda-9434-4dd7-aaad-035584b8f124"
	updatePlaceDTO.Name = "Rangel rua da Lama"
	updatePlaceDTO.Latitude = 8.826595
	updatePlaceDTO.Longitude = 13.263641
	updatePlaceDTO.Description = "Risco de inundação"
	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	place, err := placeUseCase.UpdatePlace("93247691-5c64-4c1f-a8ca-db5d76640ca9", *updatePlaceDTO)
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", place.Name)
}

func TestFindAllPlace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.Place{
		{
			ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		},
		{
			ID:          "50361691-6b99-8j2u-a8ca-db5d70912837",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		},
	}

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().FindAll().Return(data, nil)

	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	place, err := placeUseCase.FindAllPlace()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(place))
}

func TestFindPlaceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Place{
		ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
		PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:        "Rangel rua da Lama",
		Latitude:    8.825248,
		Longitude:   13.263879,
		Description: "Risco de inundação",
	}

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	place, err := placeUseCase.FindPlaceByID("93247691-5c64-4c1f-a8ca-db5d76640ca9")
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", place.Name)
}
