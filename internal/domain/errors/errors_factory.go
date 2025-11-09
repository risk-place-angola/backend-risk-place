package errors

import "errors"

// NewError creates a new contextualized error
func NewError(err error, message string) error {
	if err == nil {
		return errors.New(message)
	}
	return errors.New(message + ": " + err.Error())
}
