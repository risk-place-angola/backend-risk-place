package erfce_presenter

import (
	"github.com/labstack/echo/v4"
)

type ErfcePresenterCTX interface {
	echo.Context
}

type ErrorResponse struct {
	Message string `json:"message"`
}
