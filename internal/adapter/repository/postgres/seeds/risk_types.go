package seeds

import (
	"context"
	"database/sql"
)

func SeedRiskTypes(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO risk_types (name, description, default_radius_meters)
		VALUES
			('crime', 'Ocorrências criminais', 1000),
			('accident', 'Acidentes de trânsito', 500),
			('natural_disaster', 'Desastres naturais', 2000),
			('fire', 'Incêndios', 1500),
			('health', 'Emergências médicas', 1000),
			('infrastructure', 'Falhas de infraestrutura', 800),
			('environment', 'Riscos ambientais', 1000),
			('violence', 'Violência e agressão', 1200),
			('public_safety', 'Segurança pública', 1000),
			('traffic', 'Problemas de trânsito', 600),
			('urban_issue', 'Problemas urbanos', 500)
		ON CONFLICT (name) DO NOTHING;
	`)
	return err
}
