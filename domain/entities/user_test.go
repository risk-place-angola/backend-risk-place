package entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewUser(t *testing.T) {
	name := "any_name"
	email := "any_email"
	password := "any_password"
	user, err := entities.NewUser(name, email, password)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.NotNil(t, user, "Expected user not to be nil")

	assert.Equal(t, name, user.Name, "Expected %s, got %s", name, user.Name)
	assert.Equal(t, email, user.Email, "Expected %s, got %s", email, user.Email)
	assert.Equal(t, password, user.Password, "Expected %s, got %s", password, user.Password)
	assert.NotEmpty(t, user.ID, "Expected user ID to be non-empty")
}

func TestUserUpdate(t *testing.T) {
	name := "any_name"
	email := "any_email"
	password := "any_password"
	user, _ := entities.NewUser(name, email, password)

	newName := "newName"
	newEmail := "newEmail"
	newPassword := "newPassword"
	err := user.Update(newName, newEmail, newPassword)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.Equal(t, newName, user.Name, "Expected %s, got %s", newName, user.Name)
}
