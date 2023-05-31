package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Place struct {
	ID        string    `json:"id" valid:"uuid,required~ O ID do risco é obrigatório."`
	Latitude  float64   `json:"latitude" valid:"required~ A latitude do risco é obrigatória."`
	Longitude float64   `json:"longitude" valid:"required~ A longitude do risco é obrigatória."`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

func NewPlace(r *Place) (*Place, error) {
	place := Place{
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
	place.ID = uuid.NewV4().String()
	place.CreatedAt = time.Now()
	if err := place.isValid(); err != nil {
		return nil, err
	}
	return &place, nil
}

func (risk *Place) isValid() error {
	_, err := govalidator.ValidateStruct(risk)
	if err != nil {
		return err
	}
	return nil
}

func (risk *Place) Update(r *Place) error {
	risk.Latitude = r.Latitude
	risk.Longitude = r.Longitude
	risk.UpdatedAt = time.Now()
	return nil
}
