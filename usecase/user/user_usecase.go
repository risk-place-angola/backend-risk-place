package user

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository"
)

type UserUseCase interface {
	CreateUser(dto *CreateUserDTO) (*DTO, error)
	UpdateUser(id string, dto *UpdateUserDTO) (*DTO, error)
	FindAllUser() ([]*DTO, error)
	FindUserByID(id string) (*DTO, error)
	RemoveUser(id string) error
	Login(data *LoginDTO) (*JwtResponse, error)
	FindAllUserWarnings() ([]*DTO, error)
	FindWarningByUserID(id string) ([]*DTO, error)
}

type UserUseCaseImpl struct {
	UserRepository repository.UserRepository
}

type UserClaims struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &UserUseCaseImpl{
		UserRepository: userRepo,
	}
}

func (u *UserUseCaseImpl) CreateUser(data *CreateUserDTO) (*DTO, error) {

	user, err := entities.NewUser(data.Name, data.Phone, data.Email, data.Password)
	if err != nil {
		return nil, err
	}

	if err := u.UserRepository.Save(user); err != nil {
		return nil, err
	}

	userDto := &DTO{}

	return userDto.FromUser(user), nil
}

func (u *UserUseCaseImpl) FindAllUser() ([]*DTO, error) {
	users, err := u.UserRepository.FindAll()
	if err != nil {
		return nil, err
	}
	log.Println(users[0])
	dtoUser := &DTO{}
	dtoUsers := dtoUser.FromUserList(users)

	return dtoUsers, nil
}

func (u *UserUseCaseImpl) FindUserByID(id string) (*DTO, error) {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	dtoUser := &DTO{}

	return dtoUser.FromUser(user), nil
}

func (u *UserUseCaseImpl) UpdateUser(id string, dto *UpdateUserDTO) (*DTO, error) {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := user.Update(dto.Name, dto.Phone, dto.Email, dto.Password); err != nil {
		return nil, err
	}

	if err := u.UserRepository.Update(user); err != nil {
		return nil, err
	}

	userDTO := &DTO{}

	return userDTO.FromUser(user), nil
}

func (u *UserUseCaseImpl) RemoveUser(id string) error {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err := u.UserRepository.Delete(user.ID); err != nil {
		return err
	}

	return nil
}

func (loginUseCases *UserUseCaseImpl) Login(data *LoginDTO) (*JwtResponse, error) {

	user, err := loginUseCases.UserRepository.FindByEmail(data.Email)

	if err != nil {
		return nil, fmt.Errorf("Email ou senha incorretos")
	}

	if !user.VerifyPassword(data.Password) {
		return nil, fmt.Errorf("Email ou senha incorretos")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &UserClaims{
		ID:        user.ID,
		Email:     user.Email,
		ExpiresAt: expirationTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
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
		Name:  user.Name,
		Token: tokenString,
	}, nil
}

func (u *UserUseCaseImpl) FindAllUserWarnings() ([]*DTO, error) {
	users, err := u.UserRepository.FindAllUserWarnings()
	if err != nil {
		return nil, err
	}

	dtoUser := &DTO{}
	dtoUsers := dtoUser.FromUserList(users)

	return dtoUsers, nil
}

func (u *UserUseCaseImpl) FindWarningByUserID(id string) ([]*DTO, error) {
	users, err := u.UserRepository.FindWarningByUserID(id)
	if err != nil {
		return nil, err
	}

	dtoUser := &DTO{}
	dtoUsers := dtoUser.FromUserList(users)

	return dtoUsers, nil
}
