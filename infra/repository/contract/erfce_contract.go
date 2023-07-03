package contract

import (
	"github.com/jinzhu/gorm"
	erfceRepo "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
)

type ErfceContract interface {
	erfceContract() erfceRepo.ErfceRepository
}

type ErfceRepository struct {
	Db *gorm.DB
}

func NewErfceRepository(db *gorm.DB) *ErfceRepository {
	return &ErfceRepository{Db: db}
}

func (erfce *ErfceRepository) ErfceContract() erfceRepo.ErfceRepository {
	return repository.NewErfceRepository(erfce.Db)
}
