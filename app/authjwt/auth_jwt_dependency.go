package authjwt

import (
	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"log"
)

func AuthDependency(db *gorm.DB) IAuthService {
	authRepo := repository.NewAuthJWTRepository(db)
	authJWT := NewAuthJWT(authRepo)

	if err := authJWT.CreateCredentialJwt(); err != nil {
		log.Println(err)
	}

	return NewAuthService(authJWT)
}
