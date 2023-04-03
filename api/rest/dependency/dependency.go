package dependency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

// inject dependency

func Dependency(db *gorm.DB, echo *echo.Echo) *echo.Echo {

	DependencyRouter(db, echo)

	return echo
}
