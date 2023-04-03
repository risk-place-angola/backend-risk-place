package contract

import (
	"github.com/jinzhu/gorm"
	userRepo "github.com/risk-place-angola/backend-risk-place/domain/repository"
	"github.com/risk-place-angola/backend-risk-place/infra/repository"
)

type UserContract interface {
	UserContract() userRepo.UserRepository
}

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (user *UserRepository) UserContract() userRepo.UserRepository {
	return repository.NewUserRepository(user.Db)
}
