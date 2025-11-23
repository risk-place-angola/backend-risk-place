package model

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
	"unicode"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type AccountVerification struct {
	Code         string
	CodeVerified bool
	ExpiresAt    time.Time
	Verified     bool
}

type Address struct {
	Country      string
	Province     string
	Municipality string
	Neighborhood string
	ZipCode      string
}

type SavedLocation struct {
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
}

type User struct {
	ID                uuid.UUID
	Name              string
	Email             string
	Phone             string
	Password          string
	Latitude          float64
	Longitude         float64
	AlertRadiusMeters int
	Nif               string
	Address           Address
	HomeAddress       *SavedLocation
	WorkAddress       *SavedLocation
	DeviceToken       string
	DeviceLanguage    string
	TrustScore        int
	ReportsSubmitted  int
	ReportsVerified      int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
	AccountVerification  AccountVerification
}

const (
	DefaultAlertRadiusMeters = 1000 // in meters
	defaultCodeMax           = 900000
	defaultCodeMin           = 100000
	minimumLengthPassword    = 6
)

func NewUser(
	name,
	phone,
	email,
	password string,
) (*User, error) {
	user := &User{
		ID:                uuid.New(),
		Name:              name,
		Phone:             phone,
		Email:             email,
		Password:          password,
		AlertRadiusMeters: DefaultAlertRadiusMeters,
		CreatedAt:         time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) VerifyCode(code string) bool {
	if u.AccountVerification.Code == "" {
		return false
	}
	return u.AccountVerification.Code == code && time.Now().Before(u.AccountVerification.ExpiresAt)
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if !govalidator.IsEmail(u.Email) {
		return errors.New("email is invalid")
	}

	if u.Phone == "" {
		return errors.New("phone is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if err := u.ValidatePasswordPolicy(DefaultPasswordPolicy()); err != nil {
		return err
	}

	return nil
}

func (u *User) Update(name, phone, email string, address Address) error {
	if name != "" {
		u.Name = name
	}
	if phone != "" {
		u.Phone = phone
	}
	if email != "" {
		if !govalidator.IsEmail(email) {
			return errors.New("email is invalid")
		}
		u.Email = email
	}
	if address.Country != "" {
		u.Address.Country = address.Country
	}
	if address.Province != "" {
		u.Address.Province = address.Province
	}
	if address.Municipality != "" {
		u.Address.Municipality = address.Municipality
	}
	if address.Neighborhood != "" {
		u.Address.Neighborhood = address.Neighborhood
	}
	if address.ZipCode != "" {
		u.Address.ZipCode = address.ZipCode
	}

	u.UpdatedAt = time.Now()
	return nil
}

func GenerateConfirmationCode() string { // Generates a 6-digit code
	n, _ := rand.Int(rand.Reader, big.NewInt(defaultCodeMax))
	code := int(n.Int64()) + defaultCodeMin
	return fmt.Sprintf("%06d", code)
}

func (u *User) GenerateVerificationCode() string {
	code := GenerateConfirmationCode()
	u.AccountVerification.Code = code
	u.AccountVerification.ExpiresAt = time.Now().Add(config.CodeExpirationDuration)
	u.AccountVerification.Verified = false
	return u.AccountVerification.Code
}

func (c *User) SetGeneratedCode() error {
	genCode, err := bcrypt.GenerateFromPassword([]byte(c.AccountVerification.Code), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error generating password %s", err.Error())
	}

	c.AccountVerification.Code = string(genCode)

	return nil
}

type PasswordPolicy struct {
	MinimumLength    int
	RequireLowercase bool
	RequireNumbers   bool
	RequireSymbols   bool
	RequireUppercase bool
}

func (u *User) ValidatePasswordPolicy(passwordPolicy *PasswordPolicy) error {
	if passwordPolicy.MinimumLength > 0 && len(u.Password) < passwordPolicy.MinimumLength {
		return fmt.Errorf("password must be at least %d characters", passwordPolicy.MinimumLength)
	}

	if passwordPolicy.RequireLowercase && !ContainsLowercase(u.Password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if passwordPolicy.RequireNumbers && !ContainsNumber(u.Password) {
		return fmt.Errorf("password must contain at least one number")
	}

	if passwordPolicy.RequireSymbols && !ContainsSymbol(u.Password) {
		return fmt.Errorf("password must contain at least one symbol")
	}

	if passwordPolicy.RequireUppercase && !ContainsUppercase(u.Password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	return nil
}

func ContainsUppercase(password string) bool {
	for _, r := range password {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func ContainsSymbol(password string) bool {
	for _, r := range password {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return true
		}
	}
	return false
}

func ContainsNumber(password string) bool {
	for _, r := range password {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func ContainsLowercase(password string) bool {
	for _, r := range password {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinimumLength:    minimumLengthPassword,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSymbols:   true,
		RequireUppercase: true,
	}
}
