package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

type ErfceRepository struct {
	Db *gorm.DB
}

func NewErfceRepository(db *gorm.DB) *ErfceRepository {
	return &ErfceRepository{Db: db}
}

func (ef *ErfceRepository) Save(erfce *entities.Erfce) error {
	return ef.Db.Save(erfce).Error
}

func (ef *ErfceRepository) FindAll() ([]*entities.Erfce, error) {
	var erfce []*entities.Erfce
	err := ef.Db.Find(&erfce).Error
	return erfce, err
}

func (ef *ErfceRepository) FindByID(id string) (*entities.Erfce, error) {
	var erfce entities.Erfce
	err := ef.Db.First(&erfce, id).Error
	return &erfce, err
}

func (ef *ErfceRepository) Update(erfce *entities.Erfce) error {
	return ef.Db.Save(erfce).Error
}

func (ef *ErfceRepository) Delete(id string) error {
	return ef.Db.Delete(&entities.Erfce{}, id).Error
}

func (ef *ErfceRepository) FindByEmail(email string) (*entities.Erfce, error) {
	erfce := &entities.Erfce{}
	err := ef.Db.Where("email = ?", email).First(erfce).Error
	return erfce, err
}

func (ef *ErfceRepository) FindAllErffcesWarnings() ([]*entities.Erfce, error) {
	var erfce []*entities.Erfce
	err := ef.Db.Preload("Warnings").Find(&erfce).Error
	return erfce, err
}

func (ef *ErfceRepository) FindWarningByErfceID(id string) ([]*entities.Erfce, error) {
	var erfce []*entities.Erfce
	err := ef.Db.Preload("Warnings").Find(&erfce, "id=?", id).Error
	return erfce, err
}
