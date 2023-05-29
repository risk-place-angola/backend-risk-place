package user

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type CreateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpadateUserDTO struct {
	CreateUserDTO
}

type UserDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Token string `json:"token"`
}

func (u *UpadateUserDTO) ToUserUpdate() *entities.User {
	return &entities.User{
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *UserDTO) ToUser() *entities.User {
	return &entities.User{
		ID:       u.ID,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *UserDTO) FromUser(user *entities.User) *UserDTO {
	u.ID = user.ID
	u.Name = user.Name
	u.Phone = user.Phone
	u.Email = user.Email
	u.Password = user.Password
	u.CreatedAt = user.CreatedAt.String()
	u.UpdatedAt = user.UpdatedAt.String()
	return u
}

func (u *LoginDTO) FromUserLogin(user *entities.User) *LoginDTO {
	u.Email = user.Email
	u.Password = user.Password

	return u
}

func (u *UserDTO) FromUserList(users []*entities.User) []*UserDTO {
	var userDTO []*UserDTO
	for _, user := range users {
		userDTO = append(userDTO, u.FromUser(user))
	}
	return userDTO
}

func NewUserDTO() *UserDTO {
	return &UserDTO{}
}
