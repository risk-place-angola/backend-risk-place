package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type UserRepository interface {
	GenericRepository[entities.User]
	FindByEmail(email string) (*entities.User, error)
	FindWarningByUserID(id string) ([]*entities.User, error)
	FindAllUserWarnings() ([]*entities.User, error)
}
