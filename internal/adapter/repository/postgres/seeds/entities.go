package seeds

import (
	"context"
	"database/sql"
)

func SeedEntities(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO entities (name, entity_type, province, municipality, latitude, longitude)
	VALUES
		('ERCE Luanda', 'erce', 'Luanda', 'Luanda', -8.838333, 13.234444),
		('ERFCE Luanda', 'erfce', 'Luanda', 'Luanda', -8.838333, 13.234444);
	`)
	return err
}
