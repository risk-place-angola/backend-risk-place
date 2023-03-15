package placetype

type CreatePlaceTypeDTO struct {
	Name string `json:"name"`
}

type UpdatePlaceTypeDTO struct {
	Name string `json:"name"`
}

type PlaceTypeDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
