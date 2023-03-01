package entities_test

import (
	"testing"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewLocationType(t *testing.T) {
	locationType, err := entities.NewLocattionType("Risco")
	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	if locationType.Name != "Risco" {
		t.Error("Expected Risco, got ", locationType.Name)
	}
}

func TestLocationTypeUpdate(t *testing.T) {
	locationType, err := entities.NewLocattionType("Riscos")
	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	err = locationType.Update("Risco")
	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	if locationType.Name != "Risco" {
		t.Error("Expected Risco, got ", locationType.Name)
	}
}