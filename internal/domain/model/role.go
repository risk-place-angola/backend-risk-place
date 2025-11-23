package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID
	Name        string
	Priority    int
	Description string
}

func HighestPriorityRole(roles []Role) Role {
	role := roles[0]
	for _, r := range roles {
		if r.Priority > role.Priority {
			role = r
		}
	}
	return role
}

type UserRole struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	RoleID     uuid.UUID
	AssignedAt time.Time
}

type Permission struct {
	ID       uuid.UUID
	Resource string
	Action   string
	Code     string
}

type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	GrantedAt    time.Time
}
