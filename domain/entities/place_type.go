package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type PlaceType struct {
	ID        string    `json:"id" valid:"uuid,required~ O ID do tipo de localização é obrigatório."`
	Name      string    `json:"name" valid:"required~ O nome do tipo de localização é obrigatório."`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

func NewLocattionType(name string) (*PlaceType, error) {
	placeType := PlaceType{
		Name: name,
	}
	placeType.ID = uuid.NewV4().String()
	placeType.CreatedAt = time.Now()
	if err := placeType.isValid(); err != nil {
		return nil, err
	}
	return &placeType, nil
}

func (placeType *PlaceType) isValid() error {
	_, err := govalidator.ValidateStruct(placeType)
	if err != nil {
		return err
	}
	return nil
}

func (placeType *PlaceType) SetUpdatedAt() {
	placeType.UpdatedAt = time.Now()
}

func (placeType *PlaceType) Update(name string) error {
	placeType.Name = name
	placeType.SetUpdatedAt()
	if err := placeType.isValid(); err != nil {
		return err
	}
	return nil
}
