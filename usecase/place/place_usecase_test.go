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
	_, err := placeUseCase.CreatePlace(place_usecase.CreatePlaceDTO{
		Latitude:  8.825248,
		Longitude: 13.263879,
	})
	assert.Nil(t, err)
}

func TestUpdatePlace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Place{
		ID:        "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		Latitude:  8.825248,
		Longitude: 13.263879,
	}

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockPlaceRepository.EXPECT().Update(gomock.Any()).Return(nil)

	updatePlaceDTO := &place_usecase.UpdatePlaceDTO{}
	updatePlaceDTO.Latitude = 8.826595
	updatePlaceDTO.Longitude = 13.263641
	updatePlaceDTO.Description = "Risco de inundação"
	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	_, err := placeUseCase.UpdatePlace("93247691-5c64-4c1f-a8ca-db5d76640ca9", *updatePlaceDTO)
	assert.Nil(t, err)

}

func TestFindAllPlace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.Place{
		{
			ID:        "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			Latitude:  8.825248,
			Longitude: 13.263879,
		},
		{
			ID:        "50361691-6b99-8j2u-a8ca-db5d70912837",
			Latitude:  8.825248,
			Longitude: 13.263879,
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
		ID:        "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		Latitude:  8.825248,
		Longitude: 13.263879,
	}

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	place, err := placeUseCase.FindPlaceByID("93247691-5c64-4c1f-a8ca-db5d76640ca9")
	assert.Nil(t, err)
	assert.Equal(t, "93247691-5c64-4c1f-a8ca-db5d76640ca9", place.ID)
}
