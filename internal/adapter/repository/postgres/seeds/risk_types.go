package seeds

import (
	"context"
	"database/sql"
)

func SeedRiskTypes(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
			INSERT INTO risk_types (name, description, default_radius_meters, is_enabled)
		VALUES
			('crime', 'Ocorrências criminais', 1000, false),
			('accident', 'Acidentes de trânsito', 500, true),
			('natural_disaster', 'Desastres naturais', 2000, true),
			('fire', 'Incêndios', 1500, true),
			('health', 'Emergências médicas', 1000, true),
			('infrastructure', 'Falhas de infraestrutura', 800, true),
			('environment', 'Riscos ambientais', 1000, true),
			('violence', 'Violência e agressão', 1200, false),
			('public_safety', 'Segurança pública', 1000, true),
			('traffic', 'Problemas de trânsito', 600, true),
			('urban_issue', 'Problemas urbanos', 500, true)
		ON CONFLICT (name) DO NOTHING;
	 `)
	return err
}
