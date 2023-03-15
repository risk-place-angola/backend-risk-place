package entities_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	uuid "github.com/satori/go.uuid"
)

func TestNewRisk(t *testing.T) {
	risk, err := entities.NewRisk(&entities.Risk{
		Name:        "Viana - Estalagem",
		RiskTypeID:  uuid.NewV4().String(),
		PlaceTypeID: uuid.NewV4().String(),
		Latitude:    faker.Latitude(),
		Longitude:   faker.Longitude(),
		Description: "Homens armados assaltam a casas e estabelecimentos comerciais",
	})
	if err != nil {
		t.Errorf("Erro ao criar um novo risco: %v", err)
	}
	if risk.ID == "" {
		t.Errorf("Erro ao criar um novo risco: %v", err)
	}
}
