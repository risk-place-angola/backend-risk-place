package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type RiskType struct {
	ID          string    `json:"id" valid:"uuid,required~ O ID do tipo de risco é obrigatório."`
	Name        string    `json:"name" valid:"required~ O nome do tipo de risco é obrigatório."`
	Description string    `json:"description" valid:"required~ A descrição do tipo de risco é obrigatório."`
	CreatedAt   time.Time `json:"created_at" valid:"-"`
	UpdatedAt   time.Time `json:"updated_at" valid:"-"`
}

func NewRiskType(name, descripition string) (*RiskType, error) {
	riskType := RiskType{
		Name:        name,
		Description: descripition,
	}

	riskType.ID = uuid.NewV4().String()
	riskType.CreatedAt = time.Now()

	if err := riskType.IsValid(); err != nil {
		return nil, err
	}
	return &riskType, nil

}

func (riskType *RiskType) IsValid() error {
	_, err := govalidator.ValidateStruct(riskType)
	if err != nil {
		return err
	}
	return nil
}

func (riskType *RiskType) SetUpdatedAt() {
	riskType.UpdatedAt = time.Now()
}

func (riskType *RiskType) Update(name, descripition string) error {
	riskType.Name = name
	riskType.Description = descripition

	riskType.SetUpdatedAt()

	if err := riskType.IsValid(); err != nil {
		return err
	}
	return nil
}
