package entities

import (
	"github.com/bxcodec/faker/v3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Auth struct {
	ID        string `json:"id" valid:"uuid,required~ O ID do risco é obrigatório."`
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

func NewAuthJWTAPI() *Auth {
	return genUserAuthCreate()
}

func (auth *Auth) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password))
	return err == nil
}

func (auth *Auth) passwordEncrypt() error {
	password, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	auth.Password = string(password)
	return nil

}

func genUserAuthCreate() *Auth {
	auth := &Auth{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}
	log.Printf("Password is: %s", auth.Password)

	auth.ID = uuid.NewV4().String()
	auth.CreatedAt = time.Now()
	if err := auth.passwordEncrypt(); err != nil {
		log.Fatal(err)
	}

	return auth
}
