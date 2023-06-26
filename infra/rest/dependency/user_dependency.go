package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/router/interfaces"
	user_controller "github.com/risk-place-angola/backend-risk-place/infra/rest/user/controllers"
	user_router "github.com/risk-place-angola/backend-risk-place/infra/rest/user/router"
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
