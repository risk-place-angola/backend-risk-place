package locationtype

type CreateLocationTypeDTO struct {
	Name string `json:"name"`
}

type UpdateLocationTypeDTO struct {
	Name string `json:"name"`
}

type LocationTypeDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
