package place

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	place_controller "github.com/risk-place-angola/backend-risk-place/api/rest/place/controllers"
	place_router "github.com/risk-place-angola/backend-risk-place/api/rest/place/router"
	place_repository "github.com/risk-place-angola/backend-risk-place/infra/repository"
	place_usecase "github.com/risk-place-angola/backend-risk-place/usecase/place"
)

func PlaceDependency(db *gorm.DB, echo *echo.Echo) *echo.Echo {
	placeRepository := place_repository.NewPlaceRepository(db)
	placeUseCase := place_usecase.NewPlaceUseCase(placeRepository)
	placeController := place_controller.NewPlaceController(placeUseCase)

	placeRouterImpl := &place_router.PlaceRouterImpl{
		PlaceController: placeController,
		Echo:            echo,
	}
	placeRouter := place_router.NewPlaceRouter(placeRouterImpl)

	return placeRouter.Router()
}
