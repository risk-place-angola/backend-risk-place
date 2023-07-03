package entities_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewErfce(t *testing.T) {
	name := "risk Place"
	email := "risk@example.com"
	password := "123"
	user, err := entities.NewErfce(name, email, password)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.NotNil(t, user, "Expected user not to be nil")
	assert.Equal(t, name, user.Name, "Expected %s, got %s", name, user.Name)
	assert.Equal(t, email, user.Email, "Expected %s, got %s", email, user.Email)
	assert.True(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, "Passwords do not match")
	assert.NotEmpty(t, user.ID, "Expected user ID to be non-empty")
}

func TestErfceUpdate(t *testing.T) {
	name := "any_name"
	password := "any_password"
	user, err := entities.NewErfce(name, faker.Email(), password)
	assert.Nil(t, err, "Expected nil, got error %v", err)

	newName := "newName"
	newPassword := "newPassword"
	err = user.Update(newName, faker.Email(), newPassword)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.Equal(t, newName, user.Name, "Expected %s, got %s", newName, user.Name)
}
