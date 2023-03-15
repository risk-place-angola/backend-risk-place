package risktype_presenter

import (
	"github.com/labstack/echo/v4"
)

type RiskTypePresenterCTX interface {
	echo.Context
}

type ErrorResponse struct {
	Message string `json:"message"`
}
