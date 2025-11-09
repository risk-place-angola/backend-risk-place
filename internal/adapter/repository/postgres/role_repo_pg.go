package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	"log/slog"
)

type roleRepoPG struct {
	q sqlc.Querier
}

func (r *roleRepoPG) Create(ctx context.Context, role *model.Role) error {
	_, err := r.q.CreateRole(ctx, sqlc.CreateRoleParams{
		Name:        role.Name,
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *roleRepoPG) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	err := r.q.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *roleRepoPG) AssignRoleNameToUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	role, err := r.q.GetRoleByName(ctx, roleName)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return err
		}
		return err
	}
	err = r.q.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		UserID: userID,
		RoleID: role.ID,
	})
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			slog.Error("Role not found", slog.String("roleName", roleName), slog.Any("error", err))
			return err
		}
		slog.Error("Failed to assign role to user", slog.Any("userID", userID), slog.String("roleName", roleName), slog.Any("error", err))
		return err
	}
	return nil
}

func (r *roleRepoPG) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	roles, err := r.q.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]model.Role, 0, len(roles))
	for _, role := range roles {
		result = append(result, model.Role{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description.String,
		})
	}
	return result, nil
}

func (r *roleRepoPG) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepoPG) GetUsersByRole(ctx context.Context, roleID uuid.UUID) ([]model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepoPG) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	//TODO implement me
	panic("implement me")
}

func NewRoleRepoPG(db *sql.DB) repository.RoleRepository {
	return &roleRepoPG{
		q: sqlc.New(db),
	}
}
