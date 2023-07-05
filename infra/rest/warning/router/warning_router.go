package warning_router

import (
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/app/ws"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/middleware"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/warning/controllers"
	"os"
)

type IWaringRouter interface {
	Router() *echo.Echo
}

type WarningRouterImpl struct {
	Echo               *echo.Echo
	IWarningController warning_controllers.IWarningController
}

func NewWarningRouter(warningRouter *WarningRouterImpl) IWaringRouter {
	return &WarningRouterImpl{
		Echo:               warningRouter.Echo,
		IWarningController: warningRouter.IWarningController,
	}
}

func (w *WarningRouterImpl) Router() *echo.Echo {
	v1 := w.Echo.Group(os.Getenv("BASE_PATH"))
	{
		warning := v1.Group("/warning")
		warning.Use(middleware.AuthMiddleware())
		{
			warning.POST("", func(c echo.Context) error { return w.IWarningController.CreateWarning(c) })
			warning.GET("", func(c echo.Context) error { return w.IWarningController.FindAllWarning(c) })
			warning.PUT("/:id", func(c echo.Context) error { return w.IWarningController.UpdateWarning(c) })
			warning.GET("/:id", func(c echo.Context) error { return w.IWarningController.FindWarningByID(c) })
			warning.DELETE("/:id", func(c echo.Context) error { return w.IWarningController.RemoveWarning(c) })
			warning.GET("/ws", ws.WebsocketServer, middleware.WebsocketAuthMiddleware)
		}
	}
	return w.Echo
}
