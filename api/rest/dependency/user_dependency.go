package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/api/rest/router/interfaces"
	user_controller "github.com/risk-place-angola/backend-risk-place/api/rest/user/controllers"
	user_router "github.com/risk-place-angola/backend-risk-place/api/rest/user/router"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	user_usecase "github.com/risk-place-angola/backend-risk-place/usecase/user"
)

func UserDependency(db *gorm.DB, echo *echo.Echo) interfaces.IRouter {
	userRepository := repository.NewUserRepository(db)
	userUsecase := user_usecase.NewUserUseCase(userRepository)
	userController := user_controller.NewUserController(userUsecase)

	UserRouter := &user_router.UserRouterImpl{
		Echo:           echo,
		UserController: userController,
	}

	return user_router.NewUserRouter(UserRouter)

}
