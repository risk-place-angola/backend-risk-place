package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type RegisterUserInput struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	DeviceFCMToken string `json:"device_fcm_token,omitempty"`
	DeviceLanguage string `json:"device_language,omitempty"`
}

type RegisterUserOutput struct {
	ID uuid.UUID `json:"id"`
}

type UserProfileOutput struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	Nif         string            `json:"nif,omitempty"`
	RoleName    []string          `json:"role_name,omitempty"`
	Address     AddressDTO        `json:"address,omitempty"`
	HomeAddress *SavedLocationDTO `json:"home_address,omitempty"`
	WorkAddress *SavedLocationDTO `json:"work_address,omitempty"`
}

type UpdateUserInput struct {
	Name            string     `json:"name,omitempty"`
	Email           string     `json:"email,omitempty"`
	Phone           string     `json:"phone,omitempty"`
	Nif             string     `json:"nif,omitempty"`
	Address         AddressDTO `json:"address,omitempty"`
	CurrentPassword string     `json:"current_password,omitempty"`
	NewPassword     string     `json:"new_password,omitempty"`
}

type AddressDTO struct {
	Country      string
	Province     string
	Municipality string
	Neighborhood string
	ZipCode      string
}

type SavedLocationDTO struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateProfileRequest struct {
	HomeAddress *SavedLocationDTO `json:"home_address,omitempty"`
	WorkAddress *SavedLocationDTO `json:"work_address,omitempty"`
}

type NavigateToSavedLocationRequest struct {
	CurrentLat float64 `json:"current_lat" validate:"required,latitude"`
	CurrentLon float64 `json:"current_lon" validate:"required,longitude"`
}

type LoginInput struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	DeviceFCMToken string `json:"device_fcm_token,omitempty"`
	DeviceLanguage string `json:"device_language,omitempty"`
}

type UserSignInDTO struct {
	AccessToken         string              `json:"access_token"`
	ExpiresIn           int64               `json:"expires_in"`
	RefreshToken        string              `json:"refresh_token"`
	TokenType           string              `json:"token_type"`
	UserProfileResponse UserProfileResponse `json:"user"`
}

type UserProfileResponse struct {
	ID         string   `json:"id"`
	ActiveRole string   `json:"active_role"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	RoleName   []string `json:"role_name,omitempty"`
}

type AccessClaims struct {
	Sub        string   `json:"sub"`
	Email      string   `json:"email"`
	Roles      []string `json:"roles"`
	ActiveRole string   `json:"active_role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Sub        string       `json:"sub"`
	Roles      []model.Role `json:"roles"`
	ActiveRole string       `json:"active_role"`
	Purpose    string       `json:"purpose,omitempty"`
	jwt.RegisteredClaims
}

// ToEntityAddress converts AddressDTO to model.Address
func (a *AddressDTO) ToEntityAddress() model.Address {
	return model.Address{
		Country:      a.Country,
		Province:     a.Province,
		Municipality: a.Municipality,
		Neighborhood: a.Neighborhood,
		ZipCode:      a.ZipCode,
	}
}
