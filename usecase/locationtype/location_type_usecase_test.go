package locationtype_test

import (
	"testing"

	"github.com/golang/mock/gomock"
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
