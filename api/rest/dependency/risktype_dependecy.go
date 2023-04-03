package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	risk_type_controller "github.com/risk-place-angola/backend-risk-place/api/rest/risktype/controllers"
	risk_router "github.com/risk-place-angola/backend-risk-place/api/rest/risktype/router"
	"github.com/risk-place-angola/backend-risk-place/api/rest/router/interfaces"
	place_repository "github.com/risk-place-angola/backend-risk-place/infra/repository"
	"github.com/risk-place-angola/backend-risk-place/usecase/risktype"
)

func RiskTypeDependency(db *gorm.DB, echo *echo.Echo) interfaces.IRouter {
	riskTypeRepository := place_repository.NewRiskTypeRepository(db)
	riskTypeUseCase := risktype.NewRiskTypeUseCase(riskTypeRepository)
	riskTypeController := risk_type_controller.NewRiskTypeController(riskTypeUseCase)

	riskTypeRouterImpl := &risk_router.RiskTypeRouterImpl{
		RiskTypeController: riskTypeController,
		Echo:               echo,
	}

	return risk_router.NewRiskRouter(riskTypeRouterImpl)

}
