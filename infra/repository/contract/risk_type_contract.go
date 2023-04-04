package contract

import (
	"github.com/jinzhu/gorm"
	repositoryRiskType "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
)

type RiskTypeContract interface {
	RiskTypeContract() repositoryRiskType.RiskTypeRepository
}

type RiskTypeContractRepository struct {
	Db *gorm.DB
}

func NewRiskTypeContractRepository(db *gorm.DB) *RiskTypeContractRepository {
	return &RiskTypeContractRepository{Db: db}
}

func (l *RiskTypeContractRepository) RiskTypeContract() repositoryRiskType.RiskTypeRepository {
	return repository.NewRiskTypeRepository(l.Db)
}
