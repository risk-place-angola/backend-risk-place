package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type UserRepository interface {
	GenericRepository[entities.User]
}
