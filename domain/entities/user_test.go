package entities_test

import (
	"github.com/bxcodec/faker/v3"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

func TestNewUser(t *testing.T) {
	name := "John Doe"
	email := "johndoe@example.com"
	password := "secret123"
	phone := "912345678"
	user, err := entities.NewUser(name, phone, email, password)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.NotNil(t, user, "Expected user not to be nil")
	assert.Equal(t, name, user.Name, "Expected %s, got %s", name, user.Name)
	assert.Equal(t, email, user.Email, "Expected %s, got %s", email, user.Email)
	assert.True(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, "Passwords do not match")
	assert.NotEmpty(t, user.ID, "Expected user ID to be non-empty")
}

func TestUserUpdate(t *testing.T) {
	name := "any_name"
	password := "any_password"
	phone := "912345678"
	user, err := entities.NewUser(name, phone, faker.Email(), password)
	assert.Nil(t, err, "Expected nil, got error %v", err)

	newName := "newName"
	newPassword := "newPassword"
	err = user.Update(newName, faker.Phonenumber(), faker.Email(), newPassword)

	assert.Nil(t, err, "Expected nil, got error %v", err)
	assert.Equal(t, newName, user.Name, "Expected %s, got %s", newName, user.Name)
}
