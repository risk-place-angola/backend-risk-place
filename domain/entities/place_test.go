package entities_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewPlace(t *testing.T) {
	place, err := entities.NewPlace(&entities.Place{
		Latitude:  faker.Latitude(),
		Longitude: faker.Longitude(),
	})
	if err != nil {
		t.Errorf("Erro ao criar um novo risco: %v", err)
	}
	if place.ID == "" {
		t.Errorf("Erro ao criar um novo risco: %v", err)
	}
}
