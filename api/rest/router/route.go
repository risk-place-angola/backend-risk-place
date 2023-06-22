package router

import (
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/api/rest/middleware"
	"github.com/risk-place-angola/backend-risk-place/api/rest/router/interfaces"
	user_router "github.com/risk-place-angola/backend-risk-place/api/rest/user/router"
	warning_router "github.com/risk-place-angola/backend-risk-place/api/rest/warning/router"
	"github.com/risk-place-angola/backend-risk-place/app/ws"
	_ "github.com/risk-place-angola/backend-risk-place/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type RouterImpl struct {
	Echo *echo.Echo
	user_router.UserRouter
	warning_router.IWaringRouter
}

func NewRouter(router *RouterImpl) interfaces.IRouter {
	return &RouterImpl{
		UserRouter:    router.UserRouter,
		Echo:          router.Echo,
		IWaringRouter: router.IWaringRouter,
	}
}

func (router *RouterImpl) Router() *echo.Echo {

	router.UserRouter.Router()
	router.IWaringRouter.Router()
	router.Echo.GET("/", router.home())
	router.Echo.GET("/ws", ws.WebsocketServer, middleware.WebsocketAuthMiddleware)
	router.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	return router.Echo

}

// home is a simple handler to test our server
// @Summary Home
// @Description Home page of the API server of Risk Place Angola
// @Tags Home
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"Hello, Angola!"
// @Router / [get]
func (router *RouterImpl) home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(200, "Hello, Angola!")
	}
}
