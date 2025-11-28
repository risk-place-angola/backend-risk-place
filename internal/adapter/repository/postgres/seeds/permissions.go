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
		('user', 'read'),
		('user', 'update'),
		('user', 'delete'),
		('report', 'create'),
		('report', 'read'),
		('report', 'update'),
		('report', 'delete'),
		('report', 'verify'),
		('report', 'resolve'),
		('risk_type', 'read'),
		('risk_type', 'update'),
		('risk_type', 'manage')
	ON CONFLICT (resource, action) DO NOTHING;
	`)
	return err
}
