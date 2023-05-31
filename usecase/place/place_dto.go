package place_usecase

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type CreatePlaceDTO struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type UpdatePlaceDTO struct {
	CreatePlaceDTO
}

type PlaceDTO struct {
	ID          string  `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (u *UpdatePlaceDTO) ToPlaceUpdate() *entities.Place {
	return &entities.Place{
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
	}
}

func (r *PlaceDTO) ToPlace() *entities.Place {
	return &entities.Place{
		ID:        r.ID,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

func (r *PlaceDTO) FromPlace(place *entities.Place) *PlaceDTO {
	r.ID = place.ID
	r.Latitude = place.Latitude
	r.Longitude = place.Longitude
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
