package repository

import (
	"errors"
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
	var ReportedByUser entities.User

	err := w.Db.First(&ReportedByUser, "id=?", entity.ReportedBy).Error
	if err != nil {
		return errors.New("ReportedBy not found")
	}

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
	err := w.Db.First(&entity, "id=?", id).Error
	return &entity, err
}

func (w *WarningRepository) FindAll() ([]*entities.Warning, error) {
	var entity []*entities.Warning

	err := w.Db.Find(&entity).Error
	return entity, err
}
