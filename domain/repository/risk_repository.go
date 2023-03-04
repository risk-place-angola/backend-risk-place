package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type RiskRepository interface {
	GenericRepository[entities.Risk]
}
