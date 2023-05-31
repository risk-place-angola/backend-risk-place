package entities

import (
	"errors"
	"time"
)

type Warning struct {
	ID          string    `json:"id"`
	ReportedBy  string    `json:"reported_by"`
	User        User      `json:"user"`
	IsVictim    bool      `json:"is_victim"`
	Fact        string    `json:"fact"`
	PlaceID     string    `json:"place_id"`
	Place       Place     `json:"place"`
	IsFake      bool      `json:"is_fake"`
	IsAnonymous bool      `json:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewWarning(w *Warning) (*Warning, error) {
	warning := Warning{
		ReportedBy:  w.ReportedBy,
		IsVictim:    w.IsVictim,
		Fact:        w.Fact,
		PlaceID:     w.PlaceID,
		IsFake:      w.IsFake,
		IsAnonymous: w.IsAnonymous,
	}
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
	if warning.PlaceID == "" {
		return errors.New("o ID do local é obrigatório")
	}
	return nil
}

func (warning *Warning) SetUpdatedAt() {
	warning.UpdatedAt = time.Now()
}

func (warning *Warning) Update(reportedBy string, isVictim bool, fact string, placeID string, isFake bool, isAnonymous bool) error {
	warning.ReportedBy = reportedBy
	warning.IsVictim = isVictim
	warning.Fact = fact
	warning.PlaceID = placeID
	warning.IsFake = isFake
	warning.IsAnonymous = isAnonymous

	warning.SetUpdatedAt()
	if err := warning.isValid(); err != nil {
		return err
	}
	return nil
}
