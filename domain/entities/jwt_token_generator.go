package entities

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}

type TokenGenerator interface {
	GenerateToken(userID, email string, expiration time.Duration) (string, error)
}

func NewJwtTokenGenerator() TokenGenerator {
	return &jwtTokenGenerator{}
}

func (g *jwtTokenGenerator) GenerateToken(userID, email string, expiration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiration)

	claims := &UserClaims{
		ID:    userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %v", err)
	}

	return tokenString, nil
}
