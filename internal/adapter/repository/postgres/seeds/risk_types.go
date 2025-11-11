package seeds

import (
	"context"
	"database/sql"
)

func SeedRiskTypes(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO risk_types (name, description, default_radius_meters)
		VALUES
			('crime', 'Ocorrências criminais em geral', 1000),
			('accident', 'Acidentes de trânsito ou trabalho', 500),
			('natural_disaster', 'Desastres naturais como enchentes, deslizamentos, tempestades', 2000),
			('fire', 'Incêndios residenciais, comerciais ou florestais', 1500),
			('health', 'Emergências médicas ou surtos de doenças', 1000),
			('infrastructure', 'Falhas ou problemas em infraestrutura pública, como pontes ou energia', 1000),
			('environment', 'Riscos ambientais, poluição ou vazamento químico', 1000);
	`)
	return err
}
