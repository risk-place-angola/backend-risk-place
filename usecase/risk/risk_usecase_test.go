package risk_usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	risk_usecase "github.com/risk-place-angola/backend-risk-place/usecase/risk"
	"github.com/stretchr/testify/assert"
)

func TestCreateRisk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
	mockRiskRepository.EXPECT().Save(gomock.Any()).Return(nil)

	riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
	risk, err := riskUseCase.CreateRisk(risk_usecase.CreateRiskDTO{
		RiskTypeID:     "99bada49-09d0-4f13-b310-6f8633b38dfe",
		LocationTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:           "Rangel rua da Lama",
		Latitude:       8.825248,
		Longitude:      13.263879,
		Description:    "Risco de inundação",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", risk.Name)
}
