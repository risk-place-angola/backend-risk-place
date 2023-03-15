package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Risk struct {
	ID          string    `json:"id" valid:"uuid,required~ O ID do risco é obrigatório."`
	RiskTypeID  string    `json:"risk_type_id" valid:"uuid,required~ O ID do tipo de risco é obrigatório."`
	PlaceTypeID string    `json:"place_type_id" valid:"uuid,required~ O ID do tipo de localização é obrigatório."`
	Name        string    `json:"name" valid:"required~ O nome do risco é obrigatório."`
	Latitude    float64   `json:"latitude" valid:"required~ A latitude do risco é obrigatória."`
	Longitude   float64   `json:"longitude" valid:"required~ A longitude do risco é obrigatória."`
	Description string    `json:"description" valid:"required~ A descrição do risco é obrigatória."`
	CreatedAt   time.Time `json:"created_at" valid:"-"`
	UpdatedAt   time.Time `json:"updated_at" valid:"-"`
}

func NewRisk(r *Risk) (*Risk, error) {
	risk := Risk{
		Name:        r.Name,
		RiskTypeID:  r.RiskTypeID,
		PlaceTypeID: r.PlaceTypeID,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Description: r.Description,
	}
	risk.ID = uuid.NewV4().String()
	risk.CreatedAt = time.Now()
	if err := risk.isValid(); err != nil {
		return nil, err
	}
	return &risk, nil
}

func (risk *Risk) isValid() error {
	_, err := govalidator.ValidateStruct(risk)
	if err != nil {
		return err
	}
	return nil
}

func (risk *Risk) Update(r *Risk) error {
	risk.Name = r.Name
	risk.RiskTypeID = r.RiskTypeID
	risk.PlaceTypeID = r.PlaceTypeID
	risk.Latitude = r.Latitude
	risk.Longitude = r.Longitude
	risk.Description = r.Description
	risk.UpdatedAt = time.Now()
	return nil
}
