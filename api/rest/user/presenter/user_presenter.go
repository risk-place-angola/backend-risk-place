package user_presenter

import (
	"github.com/labstack/echo/v4"
)

type UserPresenterCTX interface {
	echo.Context
}

type ErrorResponse struct {
	Message string `json:"message"`
}
