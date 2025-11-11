package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	AssignRoleNameToUser(ctx context.Context, userID uuid.UUID, roleName string) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	GetUsersByRole(ctx context.Context, roleID uuid.UUID) ([]model.User, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}
