package risk_presenter

import (
	"github.com/labstack/echo/v4"
)

type RiskPresenterCTX interface {
	echo.Context
}

type ErrorResponse struct {
	Message string `json:"message"`
}
