package warning_usecase

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"strconv"
)

type DTO struct {
	ID          string         `json:"id"`
	ReportedBy  string         `json:"reported_by"`
	User        *entities.User `json:"user"`
	IsVictim    bool           `json:"is_victim"`
	Fact        string         `json:"fact"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	IsFake      bool           `json:"is_fake"`
	IsAnonymous bool           `json:"is_anonymous"`
}

type CreateWarningDTO struct {
	ReportedBy string `json:"reported_by"`
	IsVictim   bool   `json:"is_victim"`
	Fact       string `json:"fact"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

type UpdateWarningDTO struct {
	CreateWarningDTO
	IsFake      bool `json:"is_fake"`
	IsAnonymous bool `json:"is_anonymous"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (w *CreateWarningDTO) convertLocationStringToFloat() *Location {
	lat, _ := strconv.ParseFloat(w.Latitude, 64)
	lng, _ := strconv.ParseFloat(w.Longitude, 64)
	return &Location{
		Latitude:  lat,
		Longitude: lng,
	}
}

func (w *UpdateWarningDTO) ToWarningUpdate() *DTO {
	location := w.convertLocationStringToFloat()
	return &DTO{
		ReportedBy:  w.ReportedBy,
		IsVictim:    w.IsVictim,
		Fact:        w.Fact,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		IsFake:      w.IsFake,
		IsAnonymous: w.IsAnonymous,
	}
}

func (w *DTO) ToWarning() *DTO {
	return &DTO{
		ID:          w.ID,
		ReportedBy:  w.ReportedBy,
		IsVictim:    w.IsVictim,
		Fact:        w.Fact,
		Latitude:    w.Latitude,
		Longitude:   w.Longitude,
		IsFake:      w.IsFake,
		IsAnonymous: w.IsAnonymous,
	}
}

func (w *CreateWarningDTO) ToWarning() *DTO {
	location := w.convertLocationStringToFloat()
	return &DTO{
		ReportedBy:  w.ReportedBy,
		IsVictim:    w.IsVictim,
		Fact:        w.Fact,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		IsFake:      false,
		IsAnonymous: false,
	}
}

func (w *DTO) ToCreateWarning() *CreateWarningDTO {
	return &CreateWarningDTO{
		ReportedBy: w.ReportedBy,
		IsVictim:   w.IsVictim,
		Fact:       w.Fact,
		Latitude:   strconv.FormatFloat(w.Latitude, 'f', 6, 64),
		Longitude:  strconv.FormatFloat(w.Longitude, 'f', 6, 64),
	}
}

func (w *DTO) ToUpdateWarning() *UpdateWarningDTO {
	return &UpdateWarningDTO{
		CreateWarningDTO: CreateWarningDTO{
			ReportedBy: w.ReportedBy,
			IsVictim:   w.IsVictim,
			Fact:       w.Fact,
			Latitude:   strconv.FormatFloat(w.Latitude, 'f', 6, 64),
			Longitude:  strconv.FormatFloat(w.Longitude, 'f', 6, 64),
		},
		IsFake:      w.IsFake,
		IsAnonymous: w.IsAnonymous,
	}
}

func (w *DTO) FromWarning(warning *entities.Warning) *DTO {
	return &DTO{
		ID:          warning.ID,
		ReportedBy:  warning.ReportedBy,
		IsVictim:    warning.IsVictim,
		Fact:        warning.Fact,
		Latitude:    warning.Latitude,
		Longitude:   warning.Longitude,
		IsFake:      warning.IsFake,
		IsAnonymous: warning.IsAnonymous,
	}
}

func (w *DTO) FromWarnings(warnings []*entities.Warning) []*DTO {
	var warningsDTO []*DTO
	for _, warning := range warnings {
		warningsDTO = append(warningsDTO, w.FromWarning(warning))
	}
	return warningsDTO
}

func (w *DTO) FromCreateWarning(warning *entities.Warning) *CreateWarningDTO {
	return &CreateWarningDTO{
		ReportedBy: warning.ReportedBy,
		IsVictim:   warning.IsVictim,
		Fact:       warning.Fact,
		Latitude:   strconv.FormatFloat(warning.Latitude, 'f', 6, 64),
		Longitude:  strconv.FormatFloat(warning.Longitude, 'f', 6, 64),
	}
}

func (w *DTO) FromUpdateWarning(warning *entities.Warning) *UpdateWarningDTO {
	return &UpdateWarningDTO{
		CreateWarningDTO: CreateWarningDTO{
			ReportedBy: warning.ReportedBy,
			IsVictim:   warning.IsVictim,
			Fact:       warning.Fact,
			Latitude:   strconv.FormatFloat(warning.Latitude, 'f', 6, 64),
			Longitude:  strconv.FormatFloat(warning.Longitude, 'f', 6, 64),
		},
		IsFake:      warning.IsFake,
		IsAnonymous: warning.IsAnonymous,
	}
}
