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
	"github.com/risk-place-angola/backend-risk-place/usecase/risk/risktype"
	"github.com/stretchr/testify/assert"
)

func TestRiskTypeControllers(t *testing.T) {
	t.Run("should return 201 when create a location type", func(t *testing.T) {

		e := echo.New()

		data := entities.RiskType{
			ID:          "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:        "Assalto",
			Description: "Assalto a mão armada",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("POST", "/api/v1/risk/risktype", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		t.Log(ctx)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRiskTypeRepository := mocks.NewMockRiskTypeRepository(ctrl)
		mockRiskTypeRepository.EXPECT().Save(gomock.Any()).Return(nil)

		riskTypeUseCase := risktype.NewRiskTypeUseCase(mockRiskTypeRepository)
		riskTypeController := risk_controller.NewRiskTypeController(riskTypeUseCase)

		if assert.NoError(t, riskTypeController.RiskTypeCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}

	})
}
