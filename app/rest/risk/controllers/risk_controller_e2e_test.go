package risk_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	risk_controller "github.com/risk-place-angola/backend-risk-place/app/rest/risk/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	risk_usecase "github.com/risk-place-angola/backend-risk-place/usecase/risk"
	"github.com/stretchr/testify/assert"
)

func TestRiskController(t *testing.T) {
	t.Run("should create controller a risk", func(t *testing.T) {

		e := echo.New()

		data := &entities.Risk{
			ID:             "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:     "99bada49-09d0-4f13-b310-6f8633b38dfe",
			LocationTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:           "Rangel rua da Lama",
			Latitude:       8.825248,
			Longitude:      13.263879,
			Description:    "Risco de inundação",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("POST", "/api/v1/risk", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
		mockRiskRepository.EXPECT().Save(gomock.Any()).Return(nil)

		riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
		riskController := risk_controller.NewRiskController(riskUseCase)

		if assert.NoError(t, riskController.RiskCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}

	})

	t.Run("should find risk by id controller", func(t *testing.T) {

		e := echo.New()

		res := httptest.NewRequest("GET", "/api/v1/risk/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctx.SetPath("/api/v1/risk/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("93247691-5c64-4c1f-a8ca-db5d76640ca9")

		data := &entities.Risk{
			ID:             "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:     "99bada49-09d0-4f13-b310-6f8633b38dfe",
			LocationTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:           "Rangel rua da Lama",
			Latitude:       8.825248,
			Longitude:      13.263879,
			Description:    "Risco de inundação",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
		mockRiskRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

		riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
		riskController := risk_controller.NewRiskController(riskUseCase)

		if assert.NoError(t, riskController.RiskFindByIdController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})
}
