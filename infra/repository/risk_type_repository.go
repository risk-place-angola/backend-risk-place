package repository

import (
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/jinzhu/gorm"
)

type RiskTypeRepository struct {
	Db *gorm.DB
}

func NewRiskTypeRepository(db *gorm.DB) *RiskTypeRepository {
	return &RiskTypeRepository{Db: db}
}

func (r *RiskTypeRepository) FindAll() ([]*entities.RiskType, error) {
	var riskTypes []*entities.RiskType
	err := r.Db.Find(&riskTypes).Error
	return riskTypes, err
}

func (r *RiskTypeRepository) FindByID(id string) (*entities.RiskType, error) {
	var riskType entities.RiskType
	err := r.Db.First(&riskType, id).Error
	return &riskType, err
}

func (r *RiskTypeRepository) Save(riskType *entities.RiskType) error {
	return r.Db.Create(riskType).Error
}

func (r *RiskTypeRepository) Update(riskType *entities.RiskType) error {
	return r.Db.Save(riskType).Error
}

func (r *RiskTypeRepository) Delete(id string) error {
	return r.Db.Delete(&entities.RiskType{}, id).Error
}