package authjwt

import (
	"errors"
	"github.com/risk-place-angola/backend-risk-place/api/rest/middleware"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"log"
)

type IAuthAPI interface {
	CreateCredentialJwt() error
	Auth(username, password string) (*Token, error)
}

type AuthAPI struct {
	AuthJWTRepository *repository.AuthJWTRepository
}

type Token struct {
	Token string `json:"token"`
}

func NewAuthJWT(repoAuth *repository.AuthJWTRepository) IAuthAPI {
	return &AuthAPI{AuthJWTRepository: repoAuth}
}

func (a AuthAPI) CreateCredentialJwt() error {
	auth := entities.NewAuthJWTAPI()
	if err := a.AuthJWTRepository.FindUserIfExists(); err != nil {
		if err := a.AuthJWTRepository.Save(auth); err != nil {
			return err
		}
	}
	log.Println("Authentication User Created")
	return nil
}

func (a AuthAPI) Auth(username, password string) (*Token, error) {
	auth, err := a.AuthJWTRepository.FindByUsername(username)
	if err != nil {
		return nil, errors.New("username or password invalid")
	}
	if !auth.VerifyPassword(password) {
		return nil, errors.New("username or password invalid")
	}

	middleAuthJwt, err := middleware.NewAuthToken(username)
	if err != nil {
		return nil, err
	}

	return &Token{
		Token: middleAuthJwt,
	}, nil
}
