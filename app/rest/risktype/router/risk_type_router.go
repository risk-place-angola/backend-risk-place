package risk_router

import (
	"os"

	"github.com/labstack/echo/v4"
	risk_type_controller "github.com/risk-place-angola/backend-risk-place/app/rest/risktype/controllers"
)

type RiskTypeRouter interface {
	Router() *echo.Echo
}

type RiskTypeRouterImpl struct {
	echo               *echo.Echo
	riskTypeController risk_type_controller.RiskTypeController
}

func NewRiskRouter(riskRouter *RiskTypeRouterImpl) RiskTypeRouter {
	return &RiskTypeRouterImpl{
		echo:               riskRouter.echo,
		riskTypeController: riskRouter.riskTypeController,
	}
}

func (router *RiskTypeRouterImpl) Router() *echo.Echo {

	v1 := router.echo.Group(os.Getenv("BASE_PATH"))
	{
		riskType := v1.Group("/risktype")
		{
			riskType.POST("", func(c echo.Context) error { return router.riskTypeController.RiskTypeCreateController(c) })
			riskType.GET("", func(c echo.Context) error { return router.riskTypeController.RiskTypeFindAllController(c) })
			riskType.GET("/:id", func(c echo.Context) error { return router.riskTypeController.RiskTypeFindByIdController(c) })
			riskType.PUT("/:id", func(c echo.Context) error { return router.riskTypeController.RiskTypeUpdateController(c) })
			riskType.DELETE("/:id", func(c echo.Context) error { return router.riskTypeController.RiskTypeDeleteController(c) })
		}
	}

	return router.echo
}
