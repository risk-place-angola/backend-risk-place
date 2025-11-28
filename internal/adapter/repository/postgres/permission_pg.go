package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type PermissionPG struct {
	q sqlc.Querier
}

func (p *PermissionPG) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	perms, err := p.q.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]model.Permission, 0, len(perms))
	for _, perm := range perms {
		result = append(result, model.Permission{
			ID:       perm.ID,
			Resource: perm.Resource,
			Action:   perm.Action,
			Code:     perm.Code.String,
		})
	}
	return result, nil
}

func (p *PermissionPG) HasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	result, err := p.q.HasPermission(ctx, sqlc.HasPermissionParams{
		UserID: userID,
		Code:   sql.NullString{String: permissionCode, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

func (p *PermissionPG) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]model.Permission, error) {
	perms, err := p.q.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	result := make([]model.Permission, 0, len(perms))
	for _, perm := range perms {
		result = append(result, model.Permission{
			ID:       perm.ID,
			Resource: perm.Resource,
			Action:   perm.Action,
			Code:     perm.Code.String,
		})
	}
	return result, nil
}

func NewPermissionRepoPG(db *sql.DB) repository.PermissionRepository {
	return &PermissionPG{
		q: sqlc.New(db),
	}
}
