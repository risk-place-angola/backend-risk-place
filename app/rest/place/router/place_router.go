package place_router

import (
	"os"

	"github.com/labstack/echo/v4"
	place_controller "github.com/risk-place-angola/backend-risk-place/app/rest/place/controllers"
)

type PlaceRouter interface {
	Router() *echo.Echo
}

type PlaceRouterImpl struct {
	echo                     *echo.Echo
	placeController          place_controller.PlaceController
	placeWebsocketController *place_controller.PlaceClientManager
}

func NewPlaceRouter(placeRouter *PlaceRouterImpl) PlaceRouter {
	return &PlaceRouterImpl{
		placeController:          placeRouter.placeController,
		placeWebsocketController: placeRouter.placeWebsocketController,
		echo:                     placeRouter.echo,
	}
}

func (router *PlaceRouterImpl) Router() *echo.Echo {

	v1 := router.echo.Group(os.Getenv("BASE_PATH"))
	{
		place := v1.Group("/place")
		{
			place.POST("", func(c echo.Context) error { return router.placeController.PlaceCreateController(c) })
			place.GET("/:id", func(c echo.Context) error { return router.placeController.PlaceFindByIdController(c) })
			place.GET("/ws", func(c echo.Context) error { return router.placeWebsocketController.PlaceHandler(c) })
		}
	}

	return router.echo
}
