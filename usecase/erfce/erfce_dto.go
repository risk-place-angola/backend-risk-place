package erfce

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

type CreateErfceDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateErfceDTO struct {
	CreateErfceDTO
}

type DTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func (ef *UpdateErfceDTO) ToErfceUpdate() *entities.Erfce {
	return &entities.Erfce{
		Name:     ef.Name,
		Email:    ef.Email,
		Password: ef.Password,
	}
}

func (ef *DTO) ToErfce() *entities.Erfce {
	return &entities.Erfce{
		ID:    ef.ID,
		Name:  ef.Name,
		Email: ef.Email,
		//Password: u.Password,
	}
}

func (ef *DTO) FromErfce(erfce *entities.Erfce) *DTO {
	ef.ID = erfce.ID
	ef.Name = erfce.Name
	ef.Email = erfce.Email
	return ef
}

func (ef *LoginDTO) FromErfceLogin(erfce *entities.Erfce) *LoginDTO {
	ef.Email = erfce.Email
	ef.Password = erfce.Password

	return ef
}

func (ef *DTO) FromErfceList(erfces []*entities.Erfce) []*DTO {
	var erfceDTO []*DTO
	for _, erfce := range erfces {
		erfceDTO = append(erfceDTO, ef.FromErfce(erfce))
	}
	return erfceDTO
}

func NewErfceDTO() *DTO {
	return &DTO{}
}
