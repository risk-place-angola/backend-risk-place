package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	erfce_controller "github.com/risk-place-angola/backend-risk-place/infra/rest/erfce/controllers"
	erfce_router "github.com/risk-place-angola/backend-risk-place/infra/rest/erfce/router"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/router/interfaces"
	erfce_usecase "github.com/risk-place-angola/backend-risk-place/usecase/erfce"
)

func ErfceDependency(db *gorm.DB, echo *echo.Echo) interfaces.IRouter {
	erfceRepository := repository.NewErfceRepository(db)
	erfceUsecase := erfce_usecase.NewErfceUseCase(erfceRepository)
	erfceController := erfce_controller.NewErfceController(erfceUsecase)

	erfceRouter := &erfce_router.ErfceRouterImpl{
		Echo:            echo,
		ErfceController: erfceController,
	}

	return erfce_router.NewErfceRouter(erfceRouter)

}
