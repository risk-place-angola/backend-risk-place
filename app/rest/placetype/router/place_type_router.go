package placetype_router

import (
	"os"

	"github.com/labstack/echo/v4"
	placetype_controller "github.com/risk-place-angola/backend-risk-place/app/rest/placetype/controllers"
)

type PlaceTypeRouter interface {
	Router() *echo.Echo
}

type PlaceTypeRouterImpl struct {
	echo                *echo.Echo
	placeTypeController placetype_controller.PlaceTypeController
}

func NewPlaceTypeRouter(placeTypeController *PlaceTypeRouterImpl) PlaceTypeRouter {
	return &PlaceTypeRouterImpl{
		placeTypeController: placeTypeController.placeTypeController,
		echo:                placeTypeController.echo,
	}
}

func (router *PlaceTypeRouterImpl) Router() *echo.Echo {

	v1 := router.echo.Group(os.Getenv("BASE_PATH"))
	{
		placeType := v1.Group("/placetype")
		{
			placeType.POST("", func(c echo.Context) error { return router.placeTypeController.PlaceTypeCreateController(c) })
			placeType.PUT("/:id", func(c echo.Context) error { return router.placeTypeController.PlaceTypeUpdateController(c) })
			placeType.GET("", func(c echo.Context) error { return router.placeTypeController.PlaceTypeFindAllController(c) })
			placeType.GET("/:id", func(c echo.Context) error { return router.placeTypeController.PlaceTypeFindByIdController(c) })
			placeType.DELETE("/:id", func(c echo.Context) error { return router.placeTypeController.PlaceTypeDeleteController(c) })
		}
	}

	return router.echo
}
