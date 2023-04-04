package entities_test

import (
	"testing"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewPlaceType(t *testing.T) {
	placeType, err := entities.NewLocattionType("Risco")
	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	if placeType.Name != "Risco" {
		t.Error("Expected Risco, got ", placeType.Name)
	}
}

func TestPlaceTypeUpdate(t *testing.T) {
	placeType, err := entities.NewLocattionType("Riscos")
	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	err = placeType.Update("Risco")
	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	if placeType.Name != "Risco" {
		t.Error("Expected Risco, got ", placeType.Name)
	}
}
