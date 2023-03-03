package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
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
