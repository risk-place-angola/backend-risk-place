package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/router"
)

func DependencyRouter(db *gorm.DB, echo *echo.Echo) *echo.Echo {

	router_ := router.RouterImpl{
		UserRouter:    UserDependency(db, echo),
		Echo:          echo,
		IWaringRouter: WarningDependency(db, echo),
		ErfceRouter:   ErfceDependency(db, echo),
	}

	router.NewRouter(&router_).Router()

	return echo

}
