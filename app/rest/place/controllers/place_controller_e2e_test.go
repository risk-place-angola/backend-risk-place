package place_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	place_controller "github.com/risk-place-angola/backend-risk-place/app/rest/place/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	place_usecase "github.com/risk-place-angola/backend-risk-place/usecase/place"
	"github.com/stretchr/testify/assert"
)

func TestPlaceController(t *testing.T) {
	t.Run("should create controller a place", func(t *testing.T) {

		e := echo.New()

		data := &entities.Place{
			ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("POST", "/api/v1/place", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
		mockPlaceRepository.EXPECT().Save(gomock.Any()).Return(nil)

		placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
		placeController := place_controller.NewPlaceController(placeUseCase)

		if assert.NoError(t, placeController.PlaceCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}

	})

	t.Run("should find place by id controller", func(t *testing.T) {

		e := echo.New()

		res := httptest.NewRequest("GET", "/api/v1/place/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctx.SetPath("/api/v1/place/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("93247691-5c64-4c1f-a8ca-db5d76640ca9")

		data := &entities.Place{
			ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
		mockPlaceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

		placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
		placeController := place_controller.NewPlaceController(placeUseCase)

		if assert.NoError(t, placeController.PlaceFindByIdController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})
}
