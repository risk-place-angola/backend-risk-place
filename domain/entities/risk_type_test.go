package entities_test

import (
	"testing"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewRiskType(t *testing.T) {
	riskType, err := entities.NewRiskType("Criminalidade", "Cazenga Malueca muita criminalidade")

	if err != nil {
		t.Error("Expected nil, got ", err)
	}

	if riskType.Name != "Criminalidade" {
		t.Error("Expected Criminalidade, got ", riskType.Name)

	}

	if riskType.Description != "Cazenga Malueca muita criminalidade" {
		t.Error("Expected Cazenga Malueca muita criminalidade, got ", riskType.Description)

	}

}

func TestRiskTypeUpdate(t *testing.T) {
	riskType, err := entities.NewRiskType("Criminalidade", "Cazenga Malueca muita criminalidade")

	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	err = riskType.Update("Criminalidade", "Cazenga Malueca muita criminalidade")

	if err != nil {
		t.Error("Expected nil, got ", err)
	}

	if riskType.Name != "Criminalidade" {
		t.Error("Expected Criminalidade, got ", riskType.Name)

	}

	if riskType.Description != "Cazenga Malueca muita criminalidade" {
		t.Error("Expected Cazenga Malueca muita criminalidade, got ", riskType.Description)
	}
}
