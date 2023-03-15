package place_usecase

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type CreatePlaceDTO struct {
	RiskTypeID  string  `json:"risk_type_id"`
	PlaceTypeID string  `json:"place_type_id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type UpdatePlaceDTO struct {
	CreatePlaceDTO
}

type PlaceDTO struct {
	ID          string  `json:"id"`
	RiskTypeID  string  `json:"risk_type_id"`
	PlaceTypeID string  `json:"place_type_id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (u *UpdatePlaceDTO) ToPlaceUpdate() *entities.Place {
	return &entities.Place{
		RiskTypeID:  u.RiskTypeID,
		PlaceTypeID: u.PlaceTypeID,
		Name:        u.Name,
		Latitude:    u.Latitude,
		Longitude:   u.Longitude,
		Description: u.Description,
	}
}

func (r *PlaceDTO) ToPlace() *entities.Place {
	return &entities.Place{
		ID:          r.ID,
		RiskTypeID:  r.RiskTypeID,
		PlaceTypeID: r.PlaceTypeID,
		Name:        r.Name,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Description: r.Description,
	}
}

func (r *PlaceDTO) FromPlace(place *entities.Place) *PlaceDTO {
	r.ID = place.ID
	r.RiskTypeID = place.RiskTypeID
	r.PlaceTypeID = place.PlaceTypeID
	r.Name = place.Name
	r.Latitude = place.Latitude
	r.Longitude = place.Longitude
	r.Description = place.Description
	r.CreatedAt = place.CreatedAt.String()
	r.UpdatedAt = place.UpdatedAt.String()
	return r
}

func (r *PlaceDTO) FromPlaceList(places []*entities.Place) []*PlaceDTO {
	var placeDTOs []*PlaceDTO
	for _, place := range places {
		placeDTOs = append(placeDTOs, r.FromPlace(place))
	}
	return placeDTOs
}

func NewPlaceDTO() *PlaceDTO {
	return &PlaceDTO{}
}
