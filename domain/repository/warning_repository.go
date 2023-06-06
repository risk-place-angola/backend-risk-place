package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type IWaringRepository interface {
	GenericRepository[entities.Warning]
}
