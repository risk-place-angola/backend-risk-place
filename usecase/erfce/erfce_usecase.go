package erfce

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type ErfceUseCase interface {
	CreateErfce(dto *CreateErfceDTO) (*DTO, error)
	UpdateErfce(id string, dto *UpdateErfceDTO) (*DTO, error)
	FindAllErfce() ([]*DTO, error)
	FindErfceByID(id string) (*DTO, error)
	RemoveErfce(id string) error
	Login(data *LoginDTO) (*JwtResponse, error)
	FindAllUErfceWarnings() ([]*DTO, error)
	FindWarningByErfceID(id string) ([]*DTO, error)
}

type ErfceUseCaseImpl struct {
	ErfceRepository repository.ErfceRepository
}

type ErfceClaims struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}

func NewErfceUseCase(erfceRepo repository.ErfceRepository) ErfceUseCase {
	return &ErfceUseCaseImpl{
		ErfceRepository: erfceRepo,
	}
}

func (ef *ErfceUseCaseImpl) CreateErfce(data *CreateErfceDTO) (*DTO, error) {

	erfce, err := entities.NewErfce(data.Name, data.Email, data.Password)
	if err != nil {
		return nil, err
	}

	if err := ef.ErfceRepository.Save(erfce); err != nil {
		return nil, err
	}

	erfceDto := &DTO{}

	return erfceDto.FromErfce(erfce), nil
}

func (ef *ErfceUseCaseImpl) FindAllErfce() ([]*DTO, error) {
	erfces, err := ef.ErfceRepository.FindAll()
	if err != nil {
		return nil, err
	}
	log.Println(erfces[0])
	dtoErfce := &DTO{}
	dtoErfces := dtoErfce.FromErfceList(erfces)

	return dtoErfces, nil
}

func (ef *ErfceUseCaseImpl) FindErfceByID(id string) (*DTO, error) {
	erfce, err := ef.ErfceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	dtoErfce := &DTO{}

	return dtoErfce.FromErfce(erfce), nil
}

func (ef *ErfceUseCaseImpl) UpdateErfce(id string, dto *UpdateErfceDTO) (*DTO, error) {
	erfce, err := ef.ErfceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := erfce.Update(dto.Name, dto.Email, dto.Password); err != nil {
		return nil, err
	}

	if err := ef.ErfceRepository.Update(erfce); err != nil {
		return nil, err
	}

	erfceDTO := &DTO{}

	return erfceDTO.FromErfce(erfce), nil
}

func (ef *ErfceUseCaseImpl) RemoveErfce(id string) error {
	erfce, err := ef.ErfceRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err := ef.ErfceRepository.Delete(erfce.ID); err != nil {
		return err
	}

	return nil
}

func (loginUseCases *ErfceUseCaseImpl) Login(data *LoginDTO) (*JwtResponse, error) {

	erfce, err := loginUseCases.ErfceRepository.FindByEmail(data.Email)

	if err != nil {
		return nil, fmt.Errorf("Email ou senha incorretos")
	}

	if !erfce.VerifyPassword(data.Password) {
		return nil, fmt.Errorf("Email ou senha incorretos")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &ErfceClaims{
		ID:        erfce.ID,
		Email:     erfce.Email,
		ExpiresAt: expirationTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   erfce.ID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %v", err)
	}

	return &JwtResponse{
		Name:  erfce.Name,
		Token: tokenString,
	}, nil
}

func (ef *ErfceUseCaseImpl) FindAllUErfceWarnings() ([]*DTO, error) {
	erfces, err := ef.ErfceRepository.FindAllErffcesWarnings()
	if err != nil {
		return nil, err
	}

	dtoErfce := &DTO{}
	dtoErfces := dtoErfce.FromErfceList(erfces)

	return dtoErfces, nil
}

func (ef *ErfceUseCaseImpl) FindWarningByErfceID(id string) ([]*DTO, error) {
	erfce, err := ef.ErfceRepository.FindWarningByErfceID(id)
	if err != nil {
		return nil, err
	}

	dtoErfce := &DTO{}
	dtoErfces := dtoErfce.FromErfceList(erfce)

	return dtoErfces, nil
}
