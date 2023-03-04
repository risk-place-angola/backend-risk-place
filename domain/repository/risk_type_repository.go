package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type RiskTypeRepository interface {
	GenericRepository[entities.RiskType]
}
