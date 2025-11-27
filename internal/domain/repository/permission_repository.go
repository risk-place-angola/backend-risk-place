package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type PermissionRepository interface {
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error)
	HasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error)
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]model.Permission, error)
}
