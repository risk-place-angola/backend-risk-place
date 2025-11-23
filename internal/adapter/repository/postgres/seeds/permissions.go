package seeds

import (
	"context"
	"database/sql"
)

func SeedPermissions(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO permissions (resource, action)
	VALUES
		('user', 'create'),
		('user', 'update'),
		('user', 'delete'),
		('report', 'create'),
		('report', 'update'),
		('report', 'delete')
	ON CONFLICT (resource, action) DO NOTHING;
	`)
	return err
}
