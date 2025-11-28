-- Script para desabilitar tipos de risco sensíveis
-- Execute este script para ocultar tipos de risco sensíveis do mobile app

-- Tipos de risco que podem ser considerados sensíveis:
-- - crime: Ocorrências criminais
-- - violence: Violência e agressão
-- - fire: Incêndios (dependendo do contexto)

-- Para desabilitar um tipo de risco específico:
-- UPDATE risk_types SET is_enabled = FALSE WHERE name = 'nome_do_tipo';

-- Exemplo: Desabilitar tipos de crime e violência
-- UPDATE risk_types SET is_enabled = FALSE WHERE name IN ('crime', 'violence');

-- Para habilitar novamente:
-- UPDATE risk_types SET is_enabled = TRUE WHERE name IN ('crime', 'violence');

-- Ver todos os tipos de risco e seu status:
SELECT 
    id,
    name,
    description,
    is_enabled,
    created_at
FROM risk_types
ORDER BY name;

-- Ver quantos reports existem por tipo de risco:
SELECT 
    rt.name,
    rt.is_enabled,
    COUNT(r.id) as total_reports,
    COUNT(CASE WHEN r.status = 'pending' THEN 1 END) as pending_reports,
    COUNT(CASE WHEN r.status = 'verified' THEN 1 END) as verified_reports
FROM risk_types rt
LEFT JOIN reports r ON rt.id = r.risk_type_id
GROUP BY rt.id, rt.name, rt.is_enabled
ORDER BY total_reports DESC;

-- INSTRUÇÕES:
-- 1. Identifique quais tipos de risco devem ser desabilitados
-- 2. Execute o UPDATE para desabilitar esses tipos
-- 3. Verifique que os reports associados não aparecem mais nas queries
-- 4. Quando o problema com a Google for resolvido, habilite novamente os tipos necessários
