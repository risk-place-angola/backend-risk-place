package user_router

import (
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/middleware"
	user_controller "github.com/risk-place-angola/backend-risk-place/infra/rest/user/controllers"
	"os"
)

type UserRouter interface {
	Router() *echo.Echo
}

type UserRouterImpl struct {
	Echo           *echo.Echo
	UserController user_controller.UserController
}

func NewUserRouter(userRouter *UserRouterImpl) UserRouter {
	return &UserRouterImpl{
		Echo:           userRouter.Echo,
		UserController: userRouter.UserController,
	}
}

func (router *UserRouterImpl) Router() *echo.Echo {
	v1 := router.Echo.Group(os.Getenv("BASE_PATH"))

	{
		user := v1.Group("/user")
		{
			user.POST("", func(c echo.Context) error { return router.UserController.UserCreateController(c) })
			user.GET("", func(c echo.Context) error { return router.UserController.UserFindAllController(c) }, middleware.AuthMiddleware())
			user.GET("/:id", func(c echo.Context) error { return router.UserController.UserFindByIdController(c) }, middleware.AuthMiddleware())
			user.PUT("/:id", func(c echo.Context) error { return router.UserController.UserUpdateController(c) }, middleware.AuthMiddleware())
			user.DELETE("/:id", func(c echo.Context) error { return router.UserController.UserDeleteController(c) }, middleware.AuthMiddleware())
			user.POST("/login", func(c echo.Context) error { return router.UserController.UserLoginController(c) })
			user.GET("/warning", func(c echo.Context) error { return router.UserController.FindAllUserWarningsController(c) }, middleware.AuthMiddleware())
			user.GET("/warning/:id", func(c echo.Context) error { return router.UserController.FindWarningByUserIDController(c) }, middleware.AuthMiddleware())
		}
	}

	return router.Echo

}
