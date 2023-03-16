package contract

import (
	repositoryRiskType "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"gorm.io/gorm"
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
