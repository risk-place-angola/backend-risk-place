package seeds

import (
	"context"
	"database/sql"
)

func SeedRiskTopics(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'assalto_mao_armada', 'Assalto com arma de fogo ou branca' FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'roubo_residencia', 'Invasão e roubo em residências' FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'roubo_veiculo', 'Roubo de veículos ou carjacking' FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'furto_carteira', 'Furto de carteiras e pertences' FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'furto_telemovel', 'Furto de telemóveis' FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'vandalismo', 'Destruição de propriedade' FROM risk_types WHERE name = 'crime'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'sequestro', 'Sequestro ou rapto' FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'violencia_domestica', 'Agressão doméstica' FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'agressao_fisica', 'Agressão física ou luta' FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'tiroteio', 'Disparos de arma de fogo' FROM risk_types WHERE name = 'violence'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'acidente_viacao', 'Acidente de viação com vítimas' FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'colisao_transito', 'Colisão entre veículos' FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'atropelamento', 'Atropelamento de pedestre' FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'capotamento', 'Veículo capotado' FROM risk_types WHERE name = 'accident'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'inundacao', 'Inundação de vias ou residências' FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'deslizamento_terra', 'Deslizamento de terra ou musseque' FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'tempestade', 'Tempestade ou ventos fortes' FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'raio', 'Queda de raio' FROM risk_types WHERE name = 'natural_disaster'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'incendio_residencial', 'Incêndio em residência' FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'incendio_comercial', 'Incêndio em estabelecimento comercial' FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'incendio_mercado', 'Incêndio em mercado' FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'incendio_veiculo', 'Veículo em chamas' FROM risk_types WHERE name = 'fire'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'emergencia_medica', 'Pessoa com mal súbito' FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'surto_doenca', 'Surto de doença infecciosa' FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'acidente_trabalho', 'Acidente em obra ou local de trabalho' FROM risk_types WHERE name = 'health'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'queda_energia', 'Falta de energia elétrica' FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'queda_agua', 'Falta de água' FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'buraco_via', 'Buraco ou cratera na via' FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'semaforo_avariado', 'Semáforo avariado ou desligado' FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'cabo_solto', 'Cabo elétrico caído ou solto' FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'estrutura_risco', 'Estrutura com risco de colapso' FROM risk_types WHERE name = 'infrastructure'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'lixo_acumulado', 'Lixo acumulado na via' FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'esgoto_aberto', 'Esgoto a céu aberto' FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'poluicao_ar', 'Poluição do ar ou fumo tóxico' FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'vazamento_agua', 'Vazamento de água potável' FROM risk_types WHERE name = 'environment'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'rua_escura', 'Via sem iluminação pública' FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'zona_assalto', 'Local conhecido por assaltos frequentes' FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'vigilancia_necessaria', 'Local que necessita vigilância policial' FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'congestionamento', 'Trânsito intenso ou congestionamento' FROM risk_types WHERE name = 'traffic'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'via_bloqueada', 'Via bloqueada ou interditada' FROM risk_types WHERE name = 'traffic'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'manifestacao', 'Manifestação ou protesto' FROM risk_types WHERE name = 'traffic'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'operacao_policial', 'Operação policial em curso' FROM risk_types WHERE name = 'public_safety'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'animal_solto', 'Animal perigoso ou gado solto na via' FROM risk_types WHERE name = 'urban_issue'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	
	INSERT INTO risk_topics (risk_type_id, name, description)
	SELECT id, 'obra_sinalizacao', 'Obra sem sinalização adequada' FROM risk_types WHERE name = 'urban_issue'
	ON CONFLICT (risk_type_id, name) DO NOTHING;
	`)
	return err
}
