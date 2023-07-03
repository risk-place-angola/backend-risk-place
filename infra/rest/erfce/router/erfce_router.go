package erfce_router

import (
	"os"

	"github.com/risk-place-angola/backend-risk-place/infra/rest/middleware"

	"github.com/labstack/echo/v4"
	erfce_controller "github.com/risk-place-angola/backend-risk-place/infra/rest/erfce/controllers"
)

type ErfceRouter interface {
	Router() *echo.Echo
}

type ErfceRouterImpl struct {
	Echo            *echo.Echo
	ErfceController erfce_controller.ErfceController
}

func NewErfceRouter(erfceRouter *ErfceRouterImpl) ErfceRouter {
	return &ErfceRouterImpl{
		Echo:            erfceRouter.Echo,
		ErfceController: erfceRouter.ErfceController,
	}
}

func (router *ErfceRouterImpl) Router() *echo.Echo {
	v1 := router.Echo.Group(os.Getenv("BASE_PATH"))

	{
		user := v1.Group("/erfce")
		{
			user.POST("", func(c echo.Context) error { return router.ErfceController.ErfceCreateController(c) })
			user.GET("", func(c echo.Context) error { return router.ErfceController.ErfceFindAllController(c) }, middleware.AuthMiddleware())
			user.GET("/:id", func(c echo.Context) error { return router.ErfceController.ErfceFindByIdController(c) }, middleware.AuthMiddleware())
			user.PUT("/:id", func(c echo.Context) error { return router.ErfceController.ErfceUpdateController(c) }, middleware.AuthMiddleware())
			user.DELETE("/:id", func(c echo.Context) error { return router.ErfceController.ErfceDeleteController(c) }, middleware.AuthMiddleware())
			user.POST("/login", func(c echo.Context) error { return router.ErfceController.ErfceLoginController(c) })
			user.GET("/warning", func(c echo.Context) error { return router.ErfceController.FindAllErfceWarningsController(c) }, middleware.AuthMiddleware())
			user.GET("/warning/:id", func(c echo.Context) error { return router.ErfceController.FindWarningByErfceIDController(c) }, middleware.AuthMiddleware())
		}
	}

	return router.Echo

}
