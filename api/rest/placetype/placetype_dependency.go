package placetype

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	placetype_controller "github.com/risk-place-angola/backend-risk-place/api/rest/placetype/controllers"
	placetype_router "github.com/risk-place-angola/backend-risk-place/api/rest/placetype/router"
	place_repository "github.com/risk-place-angola/backend-risk-place/infra/repository"
	placetype_usecase "github.com/risk-place-angola/backend-risk-place/usecase/placetype"
)

func PlaceTypeDependency(db *gorm.DB, echo *echo.Echo) *echo.Echo {
	placeTypeRepository := place_repository.NewPlaceTypeRepository(db)
	placetypeUsecase := placetype_usecase.NewPlaceTypeUseCase(placeTypeRepository)
	placetypeControllers := placetype_controller.NewPlaceTypeController(placetypeUsecase)

	placetypeRouterImpl := &placetype_router.PlaceTypeRouterImpl{
		PlaceTypeController: placetypeControllers,
		Echo:                echo,
	}
	placetypeRouter := placetype_router.NewPlaceTypeRouter(placetypeRouterImpl)

	return placetypeRouter.Router()
}
