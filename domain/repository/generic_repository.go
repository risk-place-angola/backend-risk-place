package repository

type GenericRepository[T any] interface {
	Save(entity *T) error
	Update(entity *T) error
	Delete(id string) error
	FindByID(id string) (*T, error)
	FindAll() ([]*T, error)
}
