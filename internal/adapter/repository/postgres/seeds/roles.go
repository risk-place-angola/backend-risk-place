package seeds

import (
	"context"
	"database/sql"
)

func SeedRoles(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO roles (name, description)
		VALUES
			('citizen', 'Cidadão comum'),
			('erce', 'Agente da ERCE'),
			('erfce', 'Agente da ERFCE'),
			('admin', 'Administrador do sistema')
		ON CONFLICT (name) DO NOTHING;
	`)
	return err
}
