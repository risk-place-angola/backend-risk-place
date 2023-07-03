package repository

import "github.com/risk-place-angola/backend-risk-place/domain/entities"

type ErfceRepository interface {
	GenericRepository[entities.Erfce]

	FindByEmail(email string) (*entities.Erfce, error)
	FindWarningByErfceID(id string) ([]*entities.Erfce, error)
	FindAllErffcesWarnings() ([]*entities.Erfce, error)
}
