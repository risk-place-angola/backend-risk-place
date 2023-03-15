package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type PlaceRepository interface {
	GenericRepository[entities.Place]
}
