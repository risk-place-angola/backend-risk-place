package risktype_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	"github.com/risk-place-angola/backend-risk-place/usecase/risk/risktype"
	"github.com/stretchr/testify/assert"
)

func TestRiskType(t *testing.T) {
	t.Run("should create a risk type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRiskTypeRepository := mocks.NewMockRiskTypeRepository(ctrl)
		mockRiskTypeRepository.EXPECT().Save(gomock.Any()).Return(nil)

		riskTypeUseCase := risktype.NewRiskTypeUseCase(mockRiskTypeRepository)
		riskType, err := riskTypeUseCase.CreateRiskType(&risktype.CreateRiskTypeDTO{
			Name:        "Assalto",
			Description: "Assalto a mão armada",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Assalto", riskType.Name)

	})

	t.Run("should update a risk type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := &entities.RiskType{
			ID:          "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:        "Assalto",
			Description: "Assalto a mão armada",
		}

		mockRiskTypeRepository := mocks.NewMockRiskTypeRepository(ctrl)
		mockRiskTypeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockRiskTypeRepository.EXPECT().Update(gomock.Any()).Return(nil)

		riskTypeUseCase := risktype.NewRiskTypeUseCase(mockRiskTypeRepository)
		riskType, err := riskTypeUseCase.UpdateRiskType("0c1baa42-3909-4bdb-837f-a80e68232ecd", &risktype.UpdateRiskTypeDTO{
			CreateRiskTypeDTO: risktype.CreateRiskTypeDTO{
				Name:        "Assalto",
				Description: "Assalto a mão armada com arma branca",
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, "Assalto", riskType.Name)

	})
}
