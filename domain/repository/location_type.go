package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type LocationTypeRepository interface {
	GenericRepository[entities.LocationType]
}
