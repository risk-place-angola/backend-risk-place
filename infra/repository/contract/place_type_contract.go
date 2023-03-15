package contract

import (
	repositoryPlaceType "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"gorm.io/gorm"
)

type PlaceTypeContract interface {
	PlaceContract() repositoryPlaceType.PlaceTypeRepository
}

type PlaceTypeRepository struct {
	Db *gorm.DB
}

func NewPlaceTypeRepository(db *gorm.DB) *PlaceTypeRepository {
	return &PlaceTypeRepository{Db: db}
}

func (l *PlaceTypeRepository) PlaceContract() repositoryPlaceType.PlaceTypeRepository {
	return repository.NewPlaceTypeRepository(l.Db)
}
