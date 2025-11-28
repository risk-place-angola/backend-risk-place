package seeds

import (
	"context"
	"database/sql"
)

func SeedRiskTopics(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'assalto_mao_armada', 'Assalto com arma de fogo ou branca', FALSE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'roubo_residencia', 'Invasão e roubo em residências', TRUE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'roubo_veiculo', 'Roubo de veículos ou carjacking', FALSE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'furto_carteira', 'Furto de carteiras e pertences', FALSE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'furto_telemovel', 'Furto de telemóveis', FALSE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'vandalismo', 'Destruição de propriedade', FALSE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'sequestro', 'Sequestro ou rapto', TRUE FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'violencia_domestica', 'Agressão doméstica', TRUE FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'agressao_fisica', 'Agressão física ou luta', FALSE FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'tiroteio', 'Disparos de arma de fogo', FALSE FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'acidente_viacao', 'Acidente de viação com vítimas', FALSE FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'colisao_transito', 'Colisão entre veículos', FALSE FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'atropelamento', 'Atropelamento de pedestre', FALSE FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'capotamento', 'Veículo capotado', FALSE FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'inundacao', 'Inundação de vias ou residências', FALSE FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'deslizamento_terra', 'Deslizamento de terra ou musseque', FALSE FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'tempestade', 'Tempestade ou ventos fortes', FALSE FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'raio', 'Queda de raio', FALSE FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'incendio_residencial', 'Incêndio em residência', TRUE FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'incendio_comercial', 'Incêndio em estabelecimento comercial', FALSE FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'incendio_mercado', 'Incêndio em mercado', FALSE FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'incendio_veiculo', 'Veículo em chamas', FALSE FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'emergencia_medica', 'Pessoa com mal súbito', TRUE FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'surto_doenca', 'Surto de doença infecciosa', FALSE FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'acidente_trabalho', 'Acidente em obra ou local de trabalho', FALSE FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'queda_energia', 'Falta de energia elétrica', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'queda_agua', 'Falta de água', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'buraco_via', 'Buraco ou cratera na via', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'semaforo_avariado', 'Semáforo avariado ou desligado', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'cabo_solto', 'Cabo elétrico caído ou solto', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'estrutura_risco', 'Estrutura com risco de colapso', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'lixo_acumulado', 'Lixo acumulado na via', FALSE FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'esgoto_aberto', 'Esgoto a céu aberto', FALSE FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'poluicao_ar', 'Poluição do ar ou fumo tóxico', FALSE FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'vazamento_agua', 'Vazamento de água potável', FALSE FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'rua_escura', 'Via sem iluminação pública', FALSE FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'zona_assalto', 'Local conhecido por assaltos frequentes', FALSE FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'vigilancia_necessaria', 'Local que necessita vigilância policial', FALSE FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'congestionamento', 'Trânsito intenso ou congestionamento', FALSE FROM risk_types WHERE name = 'traffic'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'via_bloqueada', 'Via bloqueada ou interditada', FALSE FROM risk_types WHERE name = 'traffic'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'manifestacao', 'Manifestação ou protesto', FALSE FROM risk_types WHERE name = 'traffic'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'operacao_policial', 'Operação policial em curso', FALSE FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'animal_solto', 'Animal perigoso ou gado solto na via', FALSE FROM risk_types WHERE name = 'urban_issue'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'obra_sinalizacao', 'Obra sem sinalização adequada', FALSE FROM risk_types WHERE name = 'urban_issue'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'acidente_transito', 'Acidente de trânsito envolvendo veículos ou pedestres', FALSE FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'assalto', 'Assalto com violência, sequestro ou agressão', TRUE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'deslizamento', 'Deslizamentos de terra ou barrancos', FALSE FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'doenca_infecciosa', 'Surtos de doenças transmissíveis', FALSE FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'enchente', 'Inundações e enchentes urbanas ou rurais', FALSE FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'furtos', 'Furtos sem violência, como bolsas ou celulares', FALSE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'incendio_florestal', 'Incêndios em áreas florestais ou savanas', FALSE FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'poluicao', 'Poluição do ar, água ou solo', FALSE FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'queda', 'Quedas de pessoas em locais públicos ou privados', FALSE FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'queda_ponte', 'Desabamento ou problemas em pontes e viadutos', FALSE FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'roubo', 'Roubo em residências, comércio ou público', TRUE FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description, is_sensitive)
	SELECT id, 'vazamento_quimico', 'Vazamento de produtos químicos ou tóxicos', FALSE FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	`)
	return err
}
