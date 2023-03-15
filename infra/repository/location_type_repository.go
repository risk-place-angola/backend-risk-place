package repository

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"gorm.io/gorm"
)

type LocationTypeRepository struct {
	Db *gorm.DB
}

func NewLocationTypeRepository(db *gorm.DB) *LocationTypeRepository {
	return &LocationTypeRepository{Db: db}
}

func (r *LocationTypeRepository) FindAll() ([]*entities.LocationType, error) {
	var locationTypes []*entities.LocationType
	err := r.Db.Find(&locationTypes).Error
	return locationTypes, err
}

func (r *LocationTypeRepository) FindByID(id string) (*entities.LocationType, error) {
	var locationType entities.LocationType
	err := r.Db.First(&locationType, id).Error
	return &locationType, err
}

func (r *LocationTypeRepository) Save(locationType *entities.LocationType) error {
	return r.Db.Create(locationType).Error
}

func (r *LocationTypeRepository) Update(locationType *entities.LocationType) error {
	return r.Db.Save(locationType).Error
}

func (r *LocationTypeRepository) Delete(id string) error {
	return r.Db.Delete(&entities.LocationType{}, id).Error
}
