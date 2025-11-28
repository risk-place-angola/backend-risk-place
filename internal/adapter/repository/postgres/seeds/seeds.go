package seeds

import (
	"context"
	"database/sql"
	"fmt"
)

type Runner func(ctx context.Context, db *sql.DB) error

func RunAll(ctx context.Context, db *sql.DB) error {
	steps := []Runner{
		SeedRoles,
		SeedRiskTypes,
		SeedRiskTopics,
		SeedEntities,
		SeedPermissions,
		SeedRolePermissions,
		SeedUsers,
	}
	for i, step := range steps {
		if err := step(ctx, db); err != nil {
			return fmt.Errorf("seed step #%d failed: %w", i+1, err)
		}
	}
	return nil
}
