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
		Name: "Test",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Test", locationType.Name)
}

func TestUpdateLocationType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.LocationType{
		ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
		Name: "Test",
	}

	mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
	mockLocationTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockLocationTypeRepository.EXPECT().Update(gomock.Any()).Return(nil)

	locationTypeUseCase := locationtype.NewLocationTypeUseCase(mockLocationTypeRepository)
	locationType, err := locationTypeUseCase.UpdateLocationType("20dabe23-3541-455b-b64d-3191f2b2a303", locationtype.UpdateLocationTypeDTO{
		Name: "Test2",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Test2", locationType.Name)
}
