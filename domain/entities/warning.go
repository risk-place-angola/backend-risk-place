package entities

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Warning struct {
	ID          string         `json:"id"`
	ReportedBy  string         `json:"reported_by"`
	User        *User          `json:"user"`
	IsVictim    bool           `json:"is_victim"`
	Fact        string         `json:"fact"`
	Latitude    float64        `json:"latitude" valid:"required~ A latitude do risco é obrigatória."`
	Longitude   float64        `json:"longitude" valid:"required~ A longitude do risco é obrigatória."`
	IsFake      bool           `json:"is_fake"`
	IsAnonymous bool           `json:"is_anonymous"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index" valid:"-"`
}

func NewWarning(w *Warning) (*Warning, error) {
	warning := Warning{
		ReportedBy:  w.ReportedBy,
		IsVictim:    w.IsVictim,
		Fact:        w.Fact,
		Latitude:    w.Latitude,
		Longitude:   w.Longitude,
		IsFake:      w.IsFake,
		IsAnonymous: w.IsAnonymous,
	}
	warning.ID = uuid.NewV4().String()
	warning.CreatedAt = time.Now()
	if err := warning.isValid(); err != nil {
		return nil, err
	}
	return &warning, nil
}

func (warning *Warning) isValid() error {
	if warning.ReportedBy == "" {
		return errors.New("o ID do usuário que reportou é obrigatório")
	}
	if warning.Fact == "" {
		return errors.New("o fato é obrigatório")
	}
	if warning.Latitude == 0 {
		return errors.New("a latitude é obrigatória")
	}
	if warning.Longitude == 0 {
		return errors.New("a longitude é obrigatória")
	}
	return nil
}

func (warning *Warning) SetUpdatedAt() {
	warning.UpdatedAt = time.Now()
}
