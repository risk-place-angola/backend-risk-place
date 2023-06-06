package contract

import (
	"github.com/jinzhu/gorm"
	Irepository "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
)

type WarningContract interface {
	WarningContract() Irepository.IWaringRepository
}

type WarningContractRepository struct {
	Db *gorm.DB
}

func NewWarningContractRepository(db *gorm.DB) Irepository.IWaringRepository {
	return repository.NewWarningRepository(db)
}
