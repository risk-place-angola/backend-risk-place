package locationtype_controllers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	locationtype_controllers "github.com/risk-place-angola/backend-risk-place/app/rest/locationtype/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	locationtype_usecase "github.com/risk-place-angola/backend-risk-place/usecase/locationtype"
	"github.com/stretchr/testify/assert"
)

func TestLocationTypeController(t *testing.T) {

	t.Run("should return 201 when create a location type", func(t *testing.T) {
		e := echo.New()
		data := []byte(`{"name": "Risco", "description": "Risco de uma localização"}`)

		res := httptest.NewRequest("POST", "/api/v1/locationtype", bytes.NewBuffer(data))
		res.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
		mockLocationTypeRepository.EXPECT().Save(gomock.Any()).Return(nil)

		locationTypeUseCase := locationtype_usecase.NewLocationTypeUseCase(mockLocationTypeRepository)
		locationTypeController := locationtype_controllers.NewLocationTypeController(locationTypeUseCase)

		if assert.NoError(t, locationTypeController.LocationTypeCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}

	})

	t.Run("should return 200 when find all location types", func(t *testing.T) {
		e := echo.New()
		res := httptest.NewRequest("GET", "/api/v1/locationtype", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []*entities.LocationType{
			{
				Name: "Risco",
			},
			{
				Name: "Hospital",
			},
		}

		mockLocationTypeRepository := mocks.NewMockLocationTypeRepository(ctrl)
		mockLocationTypeRepository.EXPECT().FindAll().Return(data, nil)

		locationTypeUseCase := locationtype_usecase.NewLocationTypeUseCase(mockLocationTypeRepository)
		locationTypeController := locationtype_controllers.NewLocationTypeController(locationTypeUseCase)

		if assert.NoError(t, locationTypeController.LocationTypeFindAllController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}
	})
}