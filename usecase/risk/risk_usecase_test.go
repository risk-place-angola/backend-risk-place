package risk_usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
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
		RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
		PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:        "Rangel rua da Lama",
		Latitude:    8.825248,
		Longitude:   13.263879,
		Description: "Risco de inundação",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", risk.Name)
}

func TestUpdateRisk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Risk{
		ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
		PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:        "Rangel rua da Lama",
		Latitude:    8.825248,
		Longitude:   13.263879,
		Description: "Risco de inundação",
	}

	mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
	mockRiskRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockRiskRepository.EXPECT().Update(gomock.Any()).Return(nil)

	updateRiskDTO := &risk_usecase.UpdateRiskDTO{}
	updateRiskDTO.RiskTypeID = "99bada49-09d0-4f13-b310-6f8633b38dfe"
	updateRiskDTO.PlaceTypeID = "dd3aadda-9434-4dd7-aaad-035584b8f124"
	updateRiskDTO.Name = "Rangel rua da Lama"
	updateRiskDTO.Latitude = 8.826595
	updateRiskDTO.Longitude = 13.263641
	updateRiskDTO.Description = "Risco de inundação"
	riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
	risk, err := riskUseCase.UpdateRisk("93247691-5c64-4c1f-a8ca-db5d76640ca9", *updateRiskDTO)
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", risk.Name)
}

func TestFindAllRisk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.Risk{
		{
			ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		},
		{
			ID:          "50361691-6b99-8j2u-a8ca-db5d70912837",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		},
	}

	mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
	mockRiskRepository.EXPECT().FindAll().Return(data, nil)

	riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
	risk, err := riskUseCase.FindAllRisk()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(risk))
}

func TestFindRiskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Risk{
		ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
		PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
		Name:        "Rangel rua da Lama",
		Latitude:    8.825248,
		Longitude:   13.263879,
		Description: "Risco de inundação",
	}

	mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
	mockRiskRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
	risk, err := riskUseCase.FindRiskByID("93247691-5c64-4c1f-a8ca-db5d76640ca9")
	assert.Nil(t, err)
	assert.Equal(t, "Rangel rua da Lama", risk.Name)
}
