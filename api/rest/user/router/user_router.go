package user_router

import (
	"os"

	"github.com/labstack/echo/v4"
	user_controller "github.com/risk-place-angola/backend-risk-place/app/rest/user/controllers"
)

type UserRouter interface {
	Router() *echo.Echo
}

type UserRouterImpl struct {
	echo           *echo.Echo
	userController user_controller.UserController
}

func NewRiskRouter(userRouter *UserRouterImpl) UserRouter {
	return &UserRouterImpl{
		echo:           userRouter.echo,
		userController: userRouter.userController,
	}
}

func (router *UserRouterImpl) Router() *echo.Echo {
	v1 := router.echo.Group(os.Getenv("BASE_PATH"))

	{
		user := v1.Group("/user")
		{
			user.POST("", func(c echo.Context) error { return router.userController.UserCreateController(c) })
			user.GET("", func(c echo.Context) error { return router.userController.UserFindAllController(c) })
			user.GET("/:id", func(c echo.Context) error { return router.userController.UserFindByIdController(c) })
			user.PUT("/:id", func(c echo.Context) error { return router.userController.UserUpdateController(c) })
			user.DELETE("/:id", func(c echo.Context) error { return router.userController.UserDeleteController(c) })
			user.DELETE("/login", func(c echo.Context) error { return router.userController.UserLoginController(c) })
		}
	}

	return router.echo

}
