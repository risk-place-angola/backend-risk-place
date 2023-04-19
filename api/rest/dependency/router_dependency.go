package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/api/rest/router"
)

func DependencyRouter(db *gorm.DB, echo *echo.Echo) *echo.Echo {

	router_ := router.RouterImpl{
		PlaceTypeRouter: PlaceTypeDependency(db, echo),
		RiskTypeRouter:  RiskTypeDependency(db, echo),
		PlaceRouter:     PlaceDependency(db, echo),
		UserRouter:      UserDependency(db, echo),
		Echo:            echo,
	}

	router.NewRouter(&router_).Router()

	return echo

}
