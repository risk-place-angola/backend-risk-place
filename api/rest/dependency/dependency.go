package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/api/rest/place"
	"github.com/risk-place-angola/backend-risk-place/api/rest/placetype"
	"github.com/risk-place-angola/backend-risk-place/api/rest/risktype"
)

// inject dependency

func Dependency(db *gorm.DB, echo *echo.Echo) *echo.Echo {

	place.PlaceDependency(db, echo)
	placetype.PlaceTypeDependency(db, echo)
	risktype.RiskTypeDependency(db, echo)

	return echo
}
