package placetype_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	"github.com/risk-place-angola/backend-risk-place/usecase/placetype"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlaceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
	mockPlaceTypeRepository.EXPECT().Save(gomock.Any()).Return(nil)

	placeTypeUseCase := placetype.NewPlaceTypeUseCase(mockPlaceTypeRepository)
	placeType, err := placeTypeUseCase.CreatePlaceType(placetype.CreatePlaceTypeDTO{
		Name: "Risco",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Risco", placeType.Name)
}

func TestUpdatePlaceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.PlaceType{
		ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
		Name: "Riscos",
	}

	mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
	mockPlaceTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockPlaceTypeRepository.EXPECT().Update(gomock.Any()).Return(nil)

	placeTypeUseCase := placetype.NewPlaceTypeUseCase(mockPlaceTypeRepository)
	placeType, err := placeTypeUseCase.UpdatePlaceType("20dabe23-3541-455b-b64d-3191f2b2a303", placetype.UpdatePlaceTypeDTO{
		Name: "Risco",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Risco", placeType.Name)
}

func TestFindAllPlaceTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.PlaceType{
		{
			ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
			Name: "Risco",
		},
	}

	mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
	mockPlaceTypeRepository.EXPECT().FindAll().Return(data, nil)

	placeTypeUseCase := placetype.NewPlaceTypeUseCase(mockPlaceTypeRepository)
	placeTypes, err := placeTypeUseCase.FindAllPlaceTypes()
	assert.Nil(t, err)
	assert.Equal(t, "Risco", placeTypes[0].Name)
}

func TestFindByIdPlaceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.PlaceType{
		ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
		Name: "Risco",
	}

	mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
	mockPlaceTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	placeTypeUseCase := placetype.NewPlaceTypeUseCase(mockPlaceTypeRepository)
	placeType, err := placeTypeUseCase.FindByIdPlaceType("20dabe23-3541-455b-b64d-3191f2b2a303")
	assert.Nil(t, err)
	assert.Equal(t, "Risco", placeType.Name)
}

func TestDeletePlaceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.PlaceType{
		ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
		Name: "Risco",
	}

	mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
	mockPlaceTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockPlaceTypeRepository.EXPECT().Delete(gomock.Any()).Return(nil)

	placeTypeUseCase := placetype.NewPlaceTypeUseCase(mockPlaceTypeRepository)
	err := placeTypeUseCase.DeletePlaceType("20dabe23-3541-455b-b64d-3191f2b2a303")
	assert.Nil(t, err)
}
