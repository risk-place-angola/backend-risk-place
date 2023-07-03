package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Erfce struct {
	ID        string         `json:"id" valid:"uuid,required~ O ID é obrigatório."`
	Name      string         `json:"name" valid:"required~ O nome  é obrigatório."`
	Email     string         `json:"email" valid:"required~ O E-mail  é obrigatório."`
	Password  string         `json:"password" valid:"required~ A palavra passe  é obrigatório."`
	CreatedAt time.Time      `json:"created_at" valid:"-"`
	UpdatedAt time.Time      `json:"updated_at" valid:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" valid:"-"`
}

func NewErfce(name, email, password string) (*Erfce, error) {
	erfce := &Erfce{
		Name:     name,
		Email:    email,
		Password: password,
	}

	erfce.ID = uuid.NewV4().String()
	erfce.CreatedAt = time.Now()

	err := erfce.passwordEncrypt()
	if err != nil {
		return nil, err
	}

	if err := erfce.isValid(); err != nil {
		return nil, err
	}

	return erfce, nil
}

func (e *Erfce) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Erfce) SetUpdatedAt() {
	e.UpdatedAt = time.Now()
}

func (e *Erfce) Update(name, email, password string) error {
	e.Name = name
	e.Email = email
	e.Password = password

	err := e.passwordEncrypt()
	if err != nil {
		return err
	}

	if err := e.isValid(); err != nil {
		return err
	}

	return nil
}

func (erfce *Erfce) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(erfce.Password), []byte(password))
	return err == nil
}

func (e *Erfce) passwordEncrypt() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	e.Password = string(passwordHash)
	return nil
}
