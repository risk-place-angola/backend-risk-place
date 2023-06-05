package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/api/rest/router"
	"github.com/risk-place-angola/backend-risk-place/app/authjwt"
)

func DependencyRouter(db *gorm.DB, echo *echo.Echo) *echo.Echo {

	router_ := router.RouterImpl{
		UserRouter:   UserDependency(db, echo),
		IAuthService: authjwt.AuthDependency(db),
		Echo:         echo,
	}

	router.NewRouter(&router_).Router()

	return echo

}
