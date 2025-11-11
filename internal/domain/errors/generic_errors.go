package errors

import "errors"

// Generic / technical errors
var (
	ErrInvalidRequest       = errors.New("invalid request")
	ErrInternalServer       = errors.New("internal server error")
	ErrNotFound             = errors.New("not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrConflict             = errors.New("conflict")
	ErrBadRequest           = errors.New("bad request")
	ErrServiceUnavailable   = errors.New("application unavailable")
	ErrGatewayTimeout       = errors.New("gateway timeout")
	ErrTooManyRequests      = errors.New("too many requests")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrNotImplemented       = errors.New("not implemented")
	ErrUnprocessableEntity  = errors.New("unprocessable entity")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrInvalidInput         = errors.New("invalid input")
	ErrExpiredToken         = errors.New("token has expired")
	ErrInvalidToken         = errors.New("invalid token")
)
