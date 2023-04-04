package repository

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/jinzhu/gorm"
)

type PlaceRepository struct {
	Db *gorm.DB
}

func NewPlaceRepository(db *gorm.DB) *PlaceRepository {
	return &PlaceRepository{Db: db}
}

func (r *PlaceRepository) FindAll() ([]*entities.Place, error) {
	var places []*entities.Place
	err := r.Db.Find(&places).Error
	return places, err
}

func (r *PlaceRepository) FindByID(id string) (*entities.Place, error) {
	var place entities.Place
	err := r.Db.First(&place, id).Error
	return &place, err
}

func (r *PlaceRepository) Save(place *entities.Place) error {
	return r.Db.Create(place).Error
}

func (r *PlaceRepository) Update(place *entities.Place) error {
	return r.Db.Save(place).Error
}

func (r *PlaceRepository) Delete(id string) error {
	return r.Db.Delete(&entities.Place{}, id).Error
}
