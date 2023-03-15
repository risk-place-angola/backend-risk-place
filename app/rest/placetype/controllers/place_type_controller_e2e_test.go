package placetype_controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	placetype_controllers "github.com/risk-place-angola/backend-risk-place/app/rest/placetype/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	placetype_usecase "github.com/risk-place-angola/backend-risk-place/usecase/placetype"
	"github.com/stretchr/testify/assert"
)

func TestPlaceTypeController(t *testing.T) {

	t.Run("should return 201 when create a place type", func(t *testing.T) {
		e := echo.New()
		data := []byte(`{"name": "Risco", "description": "Risco de uma localização"}`)

		res := httptest.NewRequest("POST", "/api/v1/placetype", bytes.NewBuffer(data))
		res.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
		mockPlaceTypeRepository.EXPECT().Save(gomock.Any()).Return(nil)

		placeTypeUseCase := placetype_usecase.NewPlaceTypeUseCase(mockPlaceTypeRepository)
		placeTypeController := placetype_controllers.NewPlaceTypeController(placeTypeUseCase)

		if assert.NoError(t, placeTypeController.PlaceTypeCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}

	})

	t.Run("should return 200 when find all place types", func(t *testing.T) {
		e := echo.New()
		res := httptest.NewRequest("GET", "/api/v1/placetype", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []*entities.PlaceType{
			{
				Name: "Risco",
			},
			{
				Name: "Hospital",
			},
		}

		mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
		mockPlaceTypeRepository.EXPECT().FindAll().Return(data, nil)

		placeTypeUseCase := placetype_usecase.NewPlaceTypeUseCase(mockPlaceTypeRepository)
		placeTypeController := placetype_controllers.NewPlaceTypeController(placeTypeUseCase)

		if assert.NoError(t, placeTypeController.PlaceTypeFindAllController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}
	})

	t.Run("should return 200 when find a place type by id", func(t *testing.T) {
		e := echo.New()
		res := httptest.NewRequest("GET", "/api/v1/placetype/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("20dabe23-3541-455b-b64d-3191f2b2a303")
		ctx.SetPath("/api/v1/placetype/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := &entities.PlaceType{
			ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
			Name: "Risco",
		}

		mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
		mockPlaceTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

		placeTypeUseCase := placetype_usecase.NewPlaceTypeUseCase(mockPlaceTypeRepository)
		placeTypeController := placetype_controllers.NewPlaceTypeController(placeTypeUseCase)

		if assert.NoError(t, placeTypeController.PlaceTypeFindByIdController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}
	})

	t.Run("should return 200 when update a place type", func(t *testing.T) {
		e := echo.New()
		data := &entities.PlaceType{
			ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
			Name: "Risco",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("PUT", "/api/v1/placetype/:id", bytes.NewBuffer(jsonData))
		res.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("20dabe23-3541-455b-b64d-3191f2b2a303")
		ctx.SetPath("/api/v1/placetype/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
		mockPlaceTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockPlaceTypeRepository.EXPECT().Update(gomock.Any()).Return(nil)

		placeTypeUseCase := placetype_usecase.NewPlaceTypeUseCase(mockPlaceTypeRepository)
		placeTypeController := placetype_controllers.NewPlaceTypeController(placeTypeUseCase)

		if assert.NoError(t, placeTypeController.PlaceTypeUpdateController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}
	})

	t.Run("should return 200 when delete a place type", func(t *testing.T) {
		e := echo.New()
		res := httptest.NewRequest("DELETE", "/api/v1/placetype/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctx.SetParamNames("id")
		ctx.SetParamValues("20dabe23-3541-455b-b64d-3191f2b2a303")
		ctx.SetPath("/api/v1/placetype/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := &entities.PlaceType{
			ID:   "20dabe23-3541-455b-b64d-3191f2b2a303",
			Name: "Risco",
		}

		mockPlaceTypeRepository := mocks.NewMockPlaceTypeRepository(ctrl)
		mockPlaceTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockPlaceTypeRepository.EXPECT().Delete(gomock.Any()).Return(nil)

		placeTypeUseCase := placetype_usecase.NewPlaceTypeUseCase(mockPlaceTypeRepository)
		placeTypeController := placetype_controllers.NewPlaceTypeController(placeTypeUseCase)

		if assert.NoError(t, placeTypeController.PlaceTypeDeleteController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}
	})
}
