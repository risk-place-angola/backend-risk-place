package locationtype_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	"github.com/risk-place-angola/backend-risk-place/usecase/locationtype"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocationType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
	mockLocationTypeRepository.EXPECT().Save(gomock.Any()).Return(nil)

	locationTypeUseCase := locationtype.NewLocationTypeUseCase(mockLocationTypeRepository)
	locationType, err := locationTypeUseCase.CreateLocationType(locationtype.CreateLocationTypeDTO{
		Name: "Risco",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Risco", locationType.Name)
}

func TestUpdateLocationType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.LocationType{
		ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
		Name: "Riscos",
	}

	mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
	mockLocationTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockLocationTypeRepository.EXPECT().Update(gomock.Any()).Return(nil)

	locationTypeUseCase := locationtype.NewLocationTypeUseCase(mockLocationTypeRepository)
	locationType, err := locationTypeUseCase.UpdateLocationType("20dabe23-3541-455b-b64d-3191f2b2a303", locationtype.UpdateLocationTypeDTO{
		Name: "Risco",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Risco", locationType.Name)
}

func TestFindAllLocationTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.LocationType{
		{
			ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
			Name: "Risco",
		},
	}

	mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
	mockLocationTypeRepository.EXPECT().FindAll().Return(data, nil)

	locationTypeUseCase := locationtype.NewLocationTypeUseCase(mockLocationTypeRepository)
	locationTypes, err := locationTypeUseCase.FindAllLocationTypes()
	assert.Nil(t, err)
	assert.Equal(t, "Risco", locationTypes[0].Name)
}

func TestFindByIdLocationType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.LocationType{
		ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
		Name: "Risco",
	}

	mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
	mockLocationTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	locationTypeUseCase := locationtype.NewLocationTypeUseCase(mockLocationTypeRepository)
	locationType, err := locationTypeUseCase.FindByIdLocationType("20dabe23-3541-455b-b64d-3191f2b2a303")
	assert.Nil(t, err)
	assert.Equal(t, "Risco", locationType.Name)
}
