package warning_usecase

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"strconv"
)

type LocationConverter interface {
	convertLocationStringToFloat() *Location
}

type DTO struct {
	ID           string  `json:"id"`
	ReportedBy   string  `json:"reported_by"`
	IsVictim     bool    `json:"is_victim"`
	Fact         string  `json:"fact"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	EventState   string  `json:"event_state"`
	IsAnonymous  bool    `json:"is_anonymous"`
	StopAlerting bool    `json:"stop_alerting"`
}

type CreateWarningDTO struct {
	ID         string `json:"id"`
	ReportedBy string `json:"reported_by"`
	IsVictim   bool   `json:"is_victim"`
	Fact       string `json:"fact"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	EventState string `json:"event_state"`
}

type UpdateWarningDTO struct {
	ReportedBy   string `json:"reported_by"`
	IsVictim     bool   `json:"is_victim"`
	Fact         string `json:"fact"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	EventState   string `json:"event_state"`
	IsAnonymous  bool   `json:"is_anonymous"`
	StopAlerting bool   `json:"stop_alerting"`
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

func (w *UpdateWarningDTO) convertLocationStringToFloat() *Location {
	lat, _ := strconv.ParseFloat(w.Latitude, 64)
	lng, _ := strconv.ParseFloat(w.Longitude, 64)
	return &Location{
		Latitude:  lat,
		Longitude: lng,
	}
}

func ConvertLocationToFloat(lc LocationConverter) *Location {
	return lc.convertLocationStringToFloat()
}

func (w *UpdateWarningDTO) ToWarningUpdate() *entities.Warning {
	location := w.convertLocationStringToFloat()
	return &entities.Warning{
		ReportedBy:   w.ReportedBy,
		IsVictim:     w.IsVictim,
		Fact:         w.Fact,
		Latitude:     location.Latitude,
		Longitude:    location.Longitude,
		EventState:   entities.EventState(w.EventState),
		IsAnonymous:  w.IsAnonymous,
		StopAlerting: w.StopAlerting,
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
		EventState:  w.EventState,
		IsAnonymous: w.IsAnonymous,
	}
}

func (w *CreateWarningDTO) ToWarning() *DTO {
	location := w.convertLocationStringToFloat()
	return &DTO{
		ReportedBy:   w.ReportedBy,
		IsVictim:     w.IsVictim,
		Fact:         w.Fact,
		Latitude:     location.Latitude,
		Longitude:    location.Longitude,
		EventState:   w.EventState,
		IsAnonymous:  false,
		StopAlerting: false,
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
		ReportedBy:   w.ReportedBy,
		IsVictim:     w.IsVictim,
		Fact:         w.Fact,
		Latitude:     strconv.FormatFloat(w.Latitude, 'f', 6, 64),
		Longitude:    strconv.FormatFloat(w.Longitude, 'f', 6, 64),
		EventState:   w.EventState,
		IsAnonymous:  w.IsAnonymous,
		StopAlerting: w.StopAlerting,
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
		EventState:  string(warning.EventState),
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
		ID:         warning.ID,
		ReportedBy: warning.ReportedBy,
		IsVictim:   warning.IsVictim,
		Fact:       warning.Fact,
		Latitude:   strconv.FormatFloat(warning.Latitude, 'f', 6, 64),
		Longitude:  strconv.FormatFloat(warning.Longitude, 'f', 6, 64),
		EventState: string(warning.EventState),
	}
}

func (w *DTO) FromUpdateWarning(warning *entities.Warning) *UpdateWarningDTO {
	return &UpdateWarningDTO{
		ReportedBy:   warning.ReportedBy,
		IsVictim:     warning.IsVictim,
		Fact:         warning.Fact,
		Latitude:     strconv.FormatFloat(warning.Latitude, 'f', 6, 64),
		Longitude:    strconv.FormatFloat(warning.Longitude, 'f', 6, 64),
		EventState:   string(warning.EventState),
		IsAnonymous:  warning.IsAnonymous,
		StopAlerting: warning.StopAlerting,
	}
}
