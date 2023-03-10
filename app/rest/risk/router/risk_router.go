package risk_router

import (
	"os"

	"github.com/labstack/echo/v4"
	risk_controller "github.com/risk-place-angola/backend-risk-place/app/rest/risk/controllers"
)

type RiskRouter interface {
	Router() *echo.Echo
}

type RiskRouterImpl struct {
	echo                    *echo.Echo
	riskController          risk_controller.RiskController
	riskWebsocketController *risk_controller.RiskClientManager
	riskTypeController      risk_controller.RiskTypeController
}

func NewRiskRouter(riskRouter *RiskRouterImpl) RiskRouter {
	return &RiskRouterImpl{
		riskController:          riskRouter.riskController,
		riskWebsocketController: riskRouter.riskWebsocketController,
		echo:                    riskRouter.echo,
		riskTypeController:      riskRouter.riskTypeController,
	}
}

func (router *RiskRouterImpl) Router() *echo.Echo {

	v1 := router.echo.Group(os.Getenv("BASE_PATH"))
	{
		risk := v1.Group("/risk")
		{
			risk.POST("", func(c echo.Context) error { return router.riskController.RiskCreateController(c) })
			risk.GET("/:id", func(c echo.Context) error { return router.riskController.RiskFindByIdController(c) })
			risk.GET("/ws", func(c echo.Context) error { return router.riskWebsocketController.RiskHandler(c) })
		}
		riskType := v1.Group("risk/risktype")
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
