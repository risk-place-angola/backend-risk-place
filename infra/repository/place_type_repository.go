package repository

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"gorm.io/gorm"
)

type PlaceTypeRepository struct {
	Db *gorm.DB
}

func NewPlaceTypeRepository(db *gorm.DB) *PlaceTypeRepository {
	return &PlaceTypeRepository{Db: db}
}

func (r *PlaceTypeRepository) FindAll() ([]*entities.PlaceType, error) {
	var placeTypes []*entities.PlaceType
	err := r.Db.Find(&placeTypes).Error
	return placeTypes, err
}

func (r *PlaceTypeRepository) FindByID(id string) (*entities.PlaceType, error) {
	var placeType entities.PlaceType
	err := r.Db.First(&placeType, id).Error
	return &placeType, err
}

func (r *PlaceTypeRepository) Save(placeType *entities.PlaceType) error {
	return r.Db.Create(placeType).Error
}

func (r *PlaceTypeRepository) Update(placeType *entities.PlaceType) error {
	return r.Db.Save(placeType).Error
}

func (r *PlaceTypeRepository) Delete(id string) error {
	return r.Db.Delete(&entities.PlaceType{}, id).Error
}
