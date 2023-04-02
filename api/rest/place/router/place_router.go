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
	Echo                     *echo.Echo
	PlaceController          place_controller.PlaceController
	placeWebsocketController *place_controller.PlaceClientManager
}

func NewPlaceRouter(placeRouter *PlaceRouterImpl) PlaceRouter {
	return &PlaceRouterImpl{
		PlaceController:          placeRouter.PlaceController,
		placeWebsocketController: placeRouter.placeWebsocketController,
		Echo:                     placeRouter.Echo,
	}
}

func (router *PlaceRouterImpl) Router() *echo.Echo {

	v1 := router.Echo.Group(os.Getenv("BASE_PATH"))
	{
		place := v1.Group("/place")
		{
			place.POST("", func(c echo.Context) error { return router.PlaceController.PlaceCreateController(c) })
			place.GET("/:id", func(c echo.Context) error { return router.PlaceController.PlaceFindByIdController(c) })
			place.GET("/ws", func(c echo.Context) error { return router.placeWebsocketController.PlaceHandler(c) })
		}
	}

	return router.Echo
}
