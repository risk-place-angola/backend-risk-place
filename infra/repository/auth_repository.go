package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

type AuthJWTRepository struct {
	Db *gorm.DB
}

func NewAuthJWTRepository(db *gorm.DB) *AuthJWTRepository {
	return &AuthJWTRepository{Db: db}
}

func (a *AuthJWTRepository) Save(entity *entities.Auth) error {
	return a.Db.Save(entity).Error
}

func (a *AuthJWTRepository) FindByUsername(username string) (*entities.Auth, error) {
	var entity entities.Auth
	err := a.Db.Where("username = ?", username).First(&entity).Error
	if err != nil {
		return nil, errors.New("not found user")
	}

	return &entity, nil
}

func (a *AuthJWTRepository) FindUserIfExists() error {
	var entity entities.Auth
	var result int64
	err := a.Db.Find(&entity).Count(&result).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthJWTRepository) FindAll() ([]entities.Auth, error) {
	var entities []entities.Auth
	err := a.Db.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (a *AuthJWTRepository) DeleteAll() error {
	err := a.Db.Delete(&entities.Auth{}).Error
	if err != nil {
		return err
	}
	return nil
}
