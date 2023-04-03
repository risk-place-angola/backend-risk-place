package user_router

import (
	"os"

	"github.com/labstack/echo/v4"
	user_controller "github.com/risk-place-angola/backend-risk-place/api/rest/user/controllers"
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
			user.GET("", func(c echo.Context) error { return router.UserController.UserFindAllController(c) })
			user.GET("/:id", func(c echo.Context) error { return router.UserController.UserFindByIdController(c) })
			user.PUT("/:id", func(c echo.Context) error { return router.UserController.UserUpdateController(c) })
			user.DELETE("/:id", func(c echo.Context) error { return router.UserController.UserDeleteController(c) })
			user.DELETE("/login", func(c echo.Context) error { return router.UserController.UserLoginController(c) })
		}
	}

	return router.Echo

}
