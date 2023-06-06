package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

type WarningRepository struct {
	Db *gorm.DB
}

func NewWarningRepository(db *gorm.DB) *WarningRepository {
	return &WarningRepository{Db: db}
}

func (w *WarningRepository) Save(entity *entities.Warning) error {
	return w.Db.Create(entity).Error
}

func (w *WarningRepository) Update(entity *entities.Warning) error {
	return w.Db.Save(entity).Error
}

func (w *WarningRepository) Delete(id string) error {
	return w.Db.Delete(&entities.Warning{}, "id=?", id).Error
}

func (w *WarningRepository) FindByID(id string) (*entities.Warning, error) {
	var entity entities.Warning
	err := w.Db.First(&entity, id).Preload("User").Error
	return &entity, err
}

func (w *WarningRepository) FindAll() ([]*entities.Warning, error) {
	var entity []*entities.Warning
	err := w.Db.Preload("User").Find(&entity, "is_fake=?", false).Error
	return entity, err
}
