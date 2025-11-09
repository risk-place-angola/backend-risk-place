package seeds

import (
	"context"
	"database/sql"
	"log/slog"
)

func SeedRoles(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insere roles padrões
	_, err = tx.ExecContext(ctx, `
		WITH rl AS (
		    INSERT INTO roles (name, description)
		    VALUES
		        ('citizen', 'Cidadão comum'),
		        ('erce', 'Agente da ERCE'),
		        ('erfce', 'Agente da ERFCE'),
		        ('admin', 'Administrador do sistema')
		    RETURNING id
		)
		SELECT 1 WHERE EXISTS (SELECT 1 FROM rl)`)
	if err != nil {
		slog.Error("Failed to insert roles", slog.Any("error", err.Error()))
		return err
	}

	return tx.Commit()
}
