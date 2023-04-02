package risktype

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	risk_type_controller "github.com/risk-place-angola/backend-risk-place/app/rest/risktype/controllers"
	risk_router "github.com/risk-place-angola/backend-risk-place/app/rest/risktype/router"
	place_repository "github.com/risk-place-angola/backend-risk-place/infra/repository"
	"github.com/risk-place-angola/backend-risk-place/usecase/risktype"
)

func RiskTypeDependency(db *gorm.DB, echo *echo.Echo) *echo.Echo {
	riskTypeRepository := place_repository.NewRiskTypeRepository(db)
	riskTypeUseCase := risktype.NewRiskTypeUseCase(riskTypeRepository)
	riskTypeController := risk_type_controller.NewRiskTypeController(riskTypeUseCase)

	riskTypeRouterImpl := &risk_router.RiskTypeRouterImpl{
		RiskTypeController: riskTypeController,
		Echo:               echo,
	}
	riskTypeRouter := risk_router.NewRiskRouter(riskTypeRouterImpl)

	return riskTypeRouter.Router()

}
