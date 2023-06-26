package entities_test

import (
	"github.com/bxcodec/faker/v3"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWarning(t *testing.T) {
	name := "John Doe"
	email := "johndoe@example.com"
	password := "secret123"
	phone := "912345678"
	user, err := entities.NewUser(name, phone, email, password)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.NotNil(t, user, "Expected user not to be nil")

	ett_warning := &entities.Warning{
		ReportedBy: user.ID,
		Fact:       faker.FirstName() + ".png",
		Latitude:   faker.Latitude(),
		Longitude:  faker.Longitude(),
		EventState: "in_review",
	}

	warning, err := entities.NewWarning(ett_warning)
	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.NotNil(t, warning, "Expected warning not to be nil")
	assert.Equal(t, user.ID, warning.ReportedBy, "Expected %s, got %s", user.ID, warning.ReportedBy)
}
