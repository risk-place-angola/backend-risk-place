package risk_router

import (
	"os"

	"github.com/labstack/echo/v4"
	risk_type_controller "github.com/risk-place-angola/backend-risk-place/api/rest/risktype/controllers"
)

type RiskTypeRouter interface {
	Router() *echo.Echo
}

type RiskTypeRouterImpl struct {
	Echo               *echo.Echo
	RiskTypeController risk_type_controller.RiskTypeController
}

func NewRiskRouter(riskRouter *RiskTypeRouterImpl) RiskTypeRouter {
	return &RiskTypeRouterImpl{
		Echo:               riskRouter.Echo,
		RiskTypeController: riskRouter.RiskTypeController,
	}
}

func (router *RiskTypeRouterImpl) Router() *echo.Echo {

	v1 := router.Echo.Group(os.Getenv("BASE_PATH"))
	{
		riskType := v1.Group("/risktype")
		{
			riskType.POST("", func(c echo.Context) error { return router.RiskTypeController.RiskTypeCreateController(c) })
			riskType.GET("", func(c echo.Context) error { return router.RiskTypeController.RiskTypeFindAllController(c) })
			riskType.GET("/:id", func(c echo.Context) error { return router.RiskTypeController.RiskTypeFindByIdController(c) })
			riskType.PUT("/:id", func(c echo.Context) error { return router.RiskTypeController.RiskTypeUpdateController(c) })
			riskType.DELETE("/:id", func(c echo.Context) error { return router.RiskTypeController.RiskTypeDeleteController(c) })
		}
	}

	return router.Echo
}
