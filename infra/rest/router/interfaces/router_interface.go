package interfaces

import "github.com/labstack/echo/v4"

type IRouter interface {
	Router() *echo.Echo
}
