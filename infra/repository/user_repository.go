package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (u *UserRepository) Save(user *entities.User) error {
	return u.Db.Save(user).Error
}

func (u *UserRepository) FindAll() ([]*entities.User, error) {
	var user []*entities.User
	err := u.Db.Preload("Warnings").Find(&user).Error
	return user, err
}

func (u *UserRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	err := u.Db.Preload("Warnings").First(&user, id).Error
	return &user, err
}

func (u *UserRepository) Update(user *entities.User) error {
	return u.Db.Save(user).Error
}

func (u *UserRepository) Delete(id string) error {
	return u.Db.Delete(&entities.User{}, id).Error
}

func (u *UserRepository) FindByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	err := u.Db.Where("email = ?", email).First(user).Error
	return user, err
}

func (u *UserRepository) FindWarningByUserID(id string) ([]*entities.User, error) {
	var user []*entities.User
	err := u.Db.Preload("Warnings").Find(&user, "id=?", id).Error
	return user, err
}

func (u *UserRepository) FindAllUserWarnings() ([]*entities.User, error) {
	var user []*entities.User
	err := u.Db.Preload("Warnings").Find(&user).Error
	return user, err
}
