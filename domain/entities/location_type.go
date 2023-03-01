package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type LocationType struct {
	ID        string    `json:"id" valid:"uuid,required~ O ID do tipo de localização é obrigatório."`
	Name      string    `json:"name" valid:"required~ O nome do tipo de localização é obrigatório."`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

func NewLocattionType(name string) (*LocationType, error) {
	locationType := LocationType{
		Name: name,
	}
	locationType.ID = uuid.NewV4().String()
	locationType.CreatedAt = time.Now()
	if err := locationType.isValid(); err != nil {
		return nil, err
	}
	return &locationType, nil
}

func (locationType *LocationType) isValid() error {
	_, err := govalidator.ValidateStruct(locationType)
	if err != nil {
		return err
	}
	return nil
}

func (locationType *LocationType) SetUpdatedAt() {
	locationType.UpdatedAt = time.Now()
}

func (locationType *LocationType) Update(name string) error {
	locationType.Name = name
	locationType.SetUpdatedAt()
	if err := locationType.isValid(); err != nil {
		return err
	}
	return nil
}