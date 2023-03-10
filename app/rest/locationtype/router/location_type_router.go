package locationtype_router

import (
	"os"

	"github.com/labstack/echo/v4"
	locationtype_controller "github.com/risk-place-angola/backend-risk-place/app/rest/locationtype/controllers"
)

type LocationTypeRouter interface {
	Router() *echo.Echo
}

type LocationTypeRouterImpl struct {
	echo                   *echo.Echo
	locationTypeController locationtype_controller.LocationTypeController
}

func NewLocationTypeRouter(locationTypeController *LocationTypeRouterImpl) LocationTypeRouter {
	return &LocationTypeRouterImpl{
		locationTypeController: locationTypeController.locationTypeController,
		echo:                   locationTypeController.echo,
	}
}

func (router *LocationTypeRouterImpl) Router() *echo.Echo {

	v1 := router.echo.Group(os.Getenv("BASE_PATH"))
	{
		locationType := v1.Group("/locationtype")
		{
			locationType.POST("", func(c echo.Context) error { return router.locationTypeController.LocationTypeCreateController(c) })
			locationType.PUT("/:id", func(c echo.Context) error { return router.locationTypeController.LocationTypeUpdateController(c) })
			locationType.GET("", func(c echo.Context) error { return router.locationTypeController.LocationTypeFindAllController(c) })
			locationType.GET("/:id", func(c echo.Context) error { return router.locationTypeController.LocationTypeFindByIdController(c) })
			locationType.DELETE("/:id", func(c echo.Context) error { return router.locationTypeController.LocationTypeDeleteController(c) })
		}
	}

	return router.echo
}
