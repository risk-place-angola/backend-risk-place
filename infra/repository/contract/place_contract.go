package contract

import (
	repositoryPlace "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"github.com/jinzhu/gorm"
)

type PlaceContract interface {
	PlaceContract() repositoryPlace.PlaceRepository
}

type PlaceContractRepository struct {
	Db *gorm.DB
}

func NewPlaceContractRepository(db *gorm.DB) *PlaceContractRepository {
	return &PlaceContractRepository{Db: db}
}

func (l *PlaceContractRepository) PlaceContract() repositoryPlace.PlaceRepository {
	return repository.NewPlaceRepository(l.Db)
}
