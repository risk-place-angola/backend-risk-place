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
	Echo                *echo.Echo
	PlaceTypeController placetype_controller.PlaceTypeController
}

func NewPlaceTypeRouter(placeTypeController *PlaceTypeRouterImpl) PlaceTypeRouter {
	return &PlaceTypeRouterImpl{
		PlaceTypeController: placeTypeController.PlaceTypeController,
		Echo:                placeTypeController.Echo,
	}
}

func (router *PlaceTypeRouterImpl) Router() *echo.Echo {

	v1 := router.Echo.Group(os.Getenv("BASE_PATH"))
	{
		placeType := v1.Group("/placetype")
		{
			placeType.POST("", func(c echo.Context) error { return router.PlaceTypeController.PlaceTypeCreateController(c) })
			placeType.PUT("/:id", func(c echo.Context) error { return router.PlaceTypeController.PlaceTypeUpdateController(c) })
			placeType.GET("", func(c echo.Context) error { return router.PlaceTypeController.PlaceTypeFindAllController(c) })
			placeType.GET("/:id", func(c echo.Context) error { return router.PlaceTypeController.PlaceTypeFindByIdController(c) })
			placeType.DELETE("/:id", func(c echo.Context) error { return router.PlaceTypeController.PlaceTypeDeleteController(c) })
		}
	}

	return router.Echo
}
