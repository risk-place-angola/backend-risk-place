package seeds

import (
	"context"
	"database/sql"
)

func SeedRolePermissions(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id
	FROM roles r
	CROSS JOIN permissions p
	WHERE r.name = 'admin'
	ON CONFLICT (role_id, permission_id) DO NOTHING;

	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id
	FROM roles r
	JOIN permissions p ON (
		(p.resource = 'report' AND p.action IN ('create', 'read', 'update')) OR
		(p.resource = 'user' AND p.action IN ('read', 'update')) OR
		(p.resource = 'risk_type' AND p.action = 'read')
	)
	WHERE r.name = 'citizen'
	ON CONFLICT (role_id, permission_id) DO NOTHING;

	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id
	FROM roles r
	JOIN permissions p ON (
		(p.resource = 'report' AND p.action IN ('create', 'read', 'update', 'verify', 'resolve')) OR
		(p.resource = 'user' AND p.action IN ('read', 'update')) OR
		(p.resource = 'risk_type' AND p.action = 'read')
	)
	WHERE r.name = 'erce'
	ON CONFLICT (role_id, permission_id) DO NOTHING;

	INSERT INTO role_permissions (role_id, permission_id)
	SELECT r.id, p.id
	FROM roles r
	JOIN permissions p ON (
		(p.resource = 'report' AND p.action IN ('create', 'read', 'update', 'verify', 'resolve')) OR
		(p.resource = 'user' AND p.action IN ('read', 'update')) OR
		(p.resource = 'risk_type' AND p.action = 'read')
	)
	WHERE r.name = 'erfce'
	ON CONFLICT (role_id, permission_id) DO NOTHING;
	`)
	return err
}
