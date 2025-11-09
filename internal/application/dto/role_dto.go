package dto

import "github.com/risk-place-angola/backend-risk-place/internal/domain/model"

func ToGetRoleNames(roles []model.Role) []string {
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames
}
