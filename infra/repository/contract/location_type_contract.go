package contract

import (
	repositoryLocationType "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
	"gorm.io/gorm"
)

type LocationTypeContract interface {
	LocationContract() repositoryLocationType.LocationTypeRepository
}

type LocationTypeRepository struct {
	Db *gorm.DB
}

func NewLocationTypeRepository(db *gorm.DB) *LocationTypeRepository {
	return &LocationTypeRepository{Db: db}
}

func (l *LocationTypeRepository) LocationContract() repositoryLocationType.LocationTypeRepository {
	return repository.NewLocationTypeRepository(l.Db)
}
