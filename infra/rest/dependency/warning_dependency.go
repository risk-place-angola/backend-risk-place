package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/router/interfaces"
	warning_controllers "github.com/risk-place-angola/backend-risk-place/infra/rest/warning/controllers"
	warning_router "github.com/risk-place-angola/backend-risk-place/infra/rest/warning/router"
	warning_usecase "github.com/risk-place-angola/backend-risk-place/usecase/warning"
)

func WarningDependency(db *gorm.DB, echo *echo.Echo) interfaces.IRouter {
	warningRepository := repository.NewWarningRepository(db)
	warningUseCase := warning_usecase.NewWarningUseCase(warningRepository)
	warningController := warning_controllers.NewWarningController(warningUseCase)

	warningRouterImpl := &warning_router.WarningRouterImpl{
		IWarningController: warningController,
		Echo:               echo,
	}

	return warning_router.NewWarningRouter(warningRouterImpl)
}
