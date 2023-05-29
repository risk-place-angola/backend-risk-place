package router

import (
	"github.com/labstack/echo/v4"
	place_router "github.com/risk-place-angola/backend-risk-place/api/rest/place/router"
	placetype_router "github.com/risk-place-angola/backend-risk-place/api/rest/placetype/router"
	risk_router "github.com/risk-place-angola/backend-risk-place/api/rest/risktype/router"
	"github.com/risk-place-angola/backend-risk-place/api/rest/router/interfaces"
	user_router "github.com/risk-place-angola/backend-risk-place/api/rest/user/router"
	"github.com/risk-place-angola/backend-risk-place/app/authjwt"
	_ "github.com/risk-place-angola/backend-risk-place/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type RouterImpl struct {
	Echo *echo.Echo
	place_router.PlaceRouter
	placetype_router.PlaceTypeRouter
	risk_router.RiskTypeRouter
	user_router.UserRouter
	authjwt.IAuthService
}

func NewRouter(router *RouterImpl) interfaces.IRouter {
	return &RouterImpl{
		PlaceRouter:     router.PlaceRouter,
		PlaceTypeRouter: router.PlaceTypeRouter,
		RiskTypeRouter:  router.RiskTypeRouter,
		UserRouter:      router.UserRouter,
		Echo:            router.Echo,
		IAuthService:    router.IAuthService,
	}
}

func (router *RouterImpl) Router() *echo.Echo {

	router.PlaceRouter.Router()
	router.PlaceTypeRouter.Router()
	router.RiskTypeRouter.Router()
	router.UserRouter.Router()

	router.Echo.GET("/", router.home())
	router.Echo.GET("/auths", router.Auths)
	router.Echo.POST("/auth", router.Auth)
	router.Echo.POST("/auth/generate", router.AuthGenerateApi)
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
