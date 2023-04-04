package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type PlaceTypeRepository interface {
	GenericRepository[entities.PlaceType]
}
