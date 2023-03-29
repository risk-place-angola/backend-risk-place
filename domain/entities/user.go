package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type User struct {
	ID        string    `json:"id" valid:"uuid,required~ O ID é obrigatório."`
	Name      string    `json:"name" valid:"required~ O nome  é obrigatório."`
	Email     string    `json:"email" valid:"required~ O E-mail  é obrigatório."`
	Password  string    `json:"password" valid:"required~ A palavra passe  é obrigatório."`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

type jwtTokenGenerator struct{}

func NewUser(name, email, password string) (*User, error) {
	user := &User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()

	err := user.passwordEncrypt()
	if err != nil {
		return nil, err
	}

	if err := user.isValid(); err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) isValid() error {
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SetUpdatedAt() {
	user.UpdatedAt = time.Now()
}

func (user *User) Update(name, emai, password string) error {
	user.Name = name
	user.Email = emai
	user.Password = password

	user.SetUpdatedAt()
	if err := user.isValid(); err != nil {
		return err
	}
	return nil
}

func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (user *User) passwordEncrypt() error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(password)
	return nil

}
