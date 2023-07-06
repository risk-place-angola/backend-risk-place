package entities

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type EventState string

const (
	Pending      EventState = "pending"
	InReview     EventState = "in_review"
	Finished     EventState = "finished"
	InProgress   EventState = "in_progress"
	InResolution EventState = "in_resolution"
	Closed       EventState = "closed"
	FalseAlarm   EventState = "false_alarm"
	FalseAlert   EventState = "false_alert"
)

type NullTime struct {
	Time  time.Time `json:"time"`
	Valid bool      `json:"valid"`
}

type Warning struct {
	ID           string     `json:"id"`
	ReportedBy   string     `json:"reported_by"`
	IsVictim     bool       `json:"is_victim"`
	Fact         string     `json:"fact"`
	Latitude     float64    `json:"latitude" valid:"required~ A latitude do risco é obrigatória."`
	Longitude    float64    `json:"longitude" valid:"required~ A longitude do risco é obrigatória."`
	EventState   EventState `json:"event_state" valid:"-"`
	IsAnonymous  bool       `json:"is_anonymous"`
	StopAlerting bool       `json:"stop_alerting"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    NullTime   `json:"deleted_at" gorm:"index" valid:"-"`
}

func NewWarning(w *Warning) (*Warning, error) {
	warning := Warning{
		ReportedBy:   w.ReportedBy,
		IsVictim:     w.IsVictim,
		Fact:         w.Fact,
		Latitude:     w.Latitude,
		EventState:   w.EventState,
		Longitude:    w.Longitude,
		IsAnonymous:  w.IsAnonymous,
		StopAlerting: w.StopAlerting,
	}
	warning.ID = uuid.NewV4().String()
	warning.CreatedAt = time.Now()
	if err := warning.isValid(); err != nil {
		return nil, err
	}
	return &warning, nil
}

func (warning *Warning) isValid() error {

	if warning.EventState != Pending && warning.EventState != InReview && warning.EventState != Finished && warning.EventState != InProgress && warning.EventState != InResolution && warning.EventState != Closed && warning.EventState != FalseAlarm && warning.EventState != FalseAlert {
		return errors.New("o estado do evento é inválido")
	}

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

func (warning *Warning) Update(w *Warning) error {
	warning.ReportedBy = w.ReportedBy
	warning.IsVictim = w.IsVictim
	warning.Fact = w.Fact
	warning.Latitude = w.Latitude
	warning.Longitude = w.Longitude
	warning.EventState = w.EventState
	warning.IsAnonymous = w.IsAnonymous
	warning.StopAlerting = w.StopAlerting
	warning.SetUpdatedAt()
	return nil
}
