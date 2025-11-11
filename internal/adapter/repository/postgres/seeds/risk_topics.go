package seeds

import (
	"context"
	"database/sql"
)

func SeedRiskTopics(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	-- Crime
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'roubo', 'Roubo em residências, comércio ou público' FROM risk_types WHERE name = 'crime';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'assalto', 'Assalto com violência, sequestro ou agressão' FROM risk_types WHERE name = 'crime';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'furtos', 'Furtos sem violência, como bolsas ou celulares' FROM risk_types WHERE name = 'crime';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'vandalismo', 'Destruição de propriedade pública ou privada' FROM risk_types WHERE name = 'crime';
	
	-- Acidentes
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'acidente_transito', 'Acidente de trânsito envolvendo veículos ou pedestres' FROM risk_types WHERE name = 'accident';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'acidente_trabalho', 'Acidente em ambiente de trabalho ou obra' FROM risk_types WHERE name = 'accident';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'queda', 'Quedas de pessoas em locais públicos ou privados' FROM risk_types WHERE name = 'accident';
	
	-- Desastres Naturais
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'enchente', 'Inundações e enchentes urbanas ou rurais' FROM risk_types WHERE name = 'natural_disaster';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'deslizamento', 'Deslizamentos de terra ou barrancos' FROM risk_types WHERE name = 'natural_disaster';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'tempestade', 'Tempestades fortes, ventos e raios' FROM risk_types WHERE name = 'natural_disaster';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'incendio_florestal', 'Incêndios em áreas florestais ou savanas' FROM risk_types WHERE name = 'fire';
	
	-- Saúde
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'doenca_infecciosa', 'Surtos de doenças transmissíveis' FROM risk_types WHERE name = 'health';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'emergencia_medica', 'Situações médicas graves como parada cardíaca' FROM risk_types WHERE name = 'health';
	
	-- Infraestrutura
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'queda_ponte', 'Desabamento ou problemas em pontes e viadutos' FROM risk_types WHERE name = 'infrastructure';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'queda_energia', 'Falhas ou interrupção de fornecimento elétrico' FROM risk_types WHERE name = 'infrastructure';
	
	-- Meio ambiente
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'poluicao', 'Poluição do ar, água ou solo' FROM risk_types WHERE name = 'environment';
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'vazamento_quimico', 'Vazamento de produtos químicos ou tóxicos' FROM risk_types WHERE name = 'environment';
	`)
	return err
}
