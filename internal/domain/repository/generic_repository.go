package repository

import (
	"context"
	"github.com/google/uuid"
)

type GenericRepository[T any] interface {
	Save(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindAll(ctx context.Context) ([]*T, error)
}
