package entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewRiskType(t *testing.T) {
	name := "Criminalidade"
	description := "Cazenga Malueca muita criminalidade"
	riskType, err := entities.NewRiskType(name, description)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.NotNil(t, riskType, "Expected riskType not to be nil")

	assert.Equal(t, name, riskType.Name, "Expected %s, got %s", name, riskType.Name)
	assert.Equal(t, description, riskType.Description, "Expected %s, got %s", description, riskType.Description)
	assert.NotEmpty(t, riskType.ID, "Expected riskType ID to be non-empty")
}

func TestRiskTypeUpdate(t *testing.T) {
	name := "Criminalidade"
	description := "Cazenga Malueca muita criminalidade"
	riskType, _ := entities.NewRiskType(name, description)

	newName := "Crime"
	newDescription := "Viana muito crime"
	err := riskType.Update(newName, newDescription)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.Equal(t, newName, riskType.Name, "Expected %s, got %s", newName, riskType.Name)
}
