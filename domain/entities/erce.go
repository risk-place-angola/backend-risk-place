package entities

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Erce struct {
	ID        string         `json:"id" valid:"uuid,required~ O ID é obrigatório."`
	Name      string         `json:"name" valid:"required~ O nome  é obrigatório."`
	Email     string         `json:"email" valid:"email,required~ O email é obrigatório."`
	Password  string         `json:"password" valid:"required~ A senha é obrigatória."`
	CreatedAt time.Time      `json:"created_at" valid:"-"`
	UpdatedAt time.Time      `json:"updated_at" valid:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" valid:"-"`
}

func NewErce(name, email, password string) (*Erce, error) {
	erce := &Erce{}

	erce.Name = name
	erce.Email = email
	erce.Password = password

	erce.ID = uuid.NewV4().String()
	erce.CreatedAt = time.Now()

	err := erce.passwordEncrypt()
	if err != nil {
		return nil, err
	}

	if err := erce.isValid(); err != nil {
		return nil, err
	}

	return erce, nil
}

func (e *Erce) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Erce) SetUpdatedAt() {
	e.UpdatedAt = time.Now()
}

func (e *Erce) Update(name, email, password string) error {
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

func (e *Erce) passwordEncrypt() error {
	password, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	e.Password = string(password)
	return nil
}

func (e *Erce) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password))
	return err == nil
}
