package errors

import "errors"

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrEmailAlreadyExists     = errors.New("email already registered")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAccountNotExists   = errors.New("user account does not exist")
	ErrInvalidCode            = errors.New("invalid or unverified code")
	ErrExpiredCode            = errors.New("code has expired")
	ErrAccountNotConfirmed    = errors.New("account not confirmed, please check your phone")
	ErrAccountNotVerified     = errors.New("account not verified, please verify your account")
	ErrPersonNotFound         = errors.New("person information not found for the user")
	ErrPersonAlreadyExists    = errors.New("person information already exists for the user")
	ErrInvalidSearchQuery     = errors.New("search query is empty or too short")
	ErrAlreadyVerified        = errors.New("email already verified, no action needed")
	ErrRateLimited            = errors.New("rate limit exceeded, please try again later")
	ErrInvalidCurrentPassword = errors.New("current password is incorrect")
	ErrNoRolesAssigned        = errors.New("no roles assigned to the user")
	ErrAlertNotFound          = errors.New("alert not found")
)
