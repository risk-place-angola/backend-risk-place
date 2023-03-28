package login

import (
	"fmt"
	"time"

	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type LoginUseCase struct {
	userRepository   repository.UserRepository
	bcryptComparator entities.BcryptComparatorImpl
	TokenGenerator   entities.TokenGenerator
}

func NewLoginUseCase(userRepo repository.UserRepository, bcryptComparator entities.BcryptComparatorImpl, tokenGenerator entities.TokenGenerator) LoginUseCase {
	return LoginUseCase{
		userRepository:   userRepo,
		bcryptComparator: bcryptComparator,
		TokenGenerator:   tokenGenerator,
	}
}

func (loginUseCase *LoginUseCase) Login(data *LoginDTO) (string, error) {
	user, err := loginUseCase.userRepository.FindByEmail(data.Email)

	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := loginUseCase.bcryptComparator.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := loginUseCase.TokenGenerator.GenerateToken(user.ID, user.Email, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %v", err)
	}

	return token, nil
}
