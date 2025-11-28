-- Script de Debug: Anonymous User Tracking
-- Use este script para diagnosticar problemas de visibilidade entre usuários anônimos

-- ============================================================================
-- 1. VERIFICAR LOCALIZAÇÕES ANÔNIMAS RECENTES
-- ============================================================================
-- Mostra todas as localizações de usuários anônimos nos últimos 5 minutos
SELECT 
    user_id,
    device_id,
    latitude,
    longitude,
    is_anonymous,
    last_update,
    NOW() - last_update as idade,
    avatar_id,
    color,
    speed,
    heading
FROM user_locations
WHERE is_anonymous = true
  AND last_update > NOW() - INTERVAL '5 minutes'
ORDER BY last_update DESC;

-- ============================================================================
-- 2. CONTAR USUÁRIOS ATIVOS (ÚLTIMOS 30 SEGUNDOS)
-- ============================================================================
SELECT 
    is_anonymous,
    COUNT(*) as total_usuarios
FROM user_locations
WHERE last_update > NOW() - INTERVAL '30 seconds'
GROUP BY is_anonymous;

-- ============================================================================
-- 3. VERIFICAR DISTÂNCIA ENTRE DOIS USUÁRIOS ESPECÍFICOS
-- ============================================================================
-- Substitua os UUIDs pelos device_ids que você está testando:
WITH user_a AS (
    SELECT user_id, latitude, longitude, location
    FROM user_locations
    WHERE user_id = '8bc22b4a-ad8a-4365-950c-bfd5fc7ec744'  -- Mobile A
),
user_b AS (
    SELECT user_id, latitude, longitude, location
    FROM user_locations
    WHERE user_id = 'c952b59e-dc44-4ec5-a944-2d8323b6ba5a'  -- Mobile B
)
SELECT 
    user_a.user_id as user_a_id,
    user_b.user_id as user_b_id,
    user_a.latitude as lat_a,
    user_a.longitude as lon_a,
    user_b.latitude as lat_b,
    user_b.longitude as lon_b,
    ST_Distance(user_a.location, user_b.location) as distancia_metros,
    CASE 
        WHEN ST_Distance(user_a.location, user_b.location) < 5000 
        THEN '✅ DENTRO DO RAIO (5km)'
        ELSE '❌ FORA DO RAIO'
    END as status
FROM user_a, user_b;

-- ============================================================================
-- 4. SIMULAR QUERY DE NEARBY USERS PARA USER A
-- ============================================================================
-- Esta é a MESMA query que o backend executa
-- Substitua o UUID pelo device_id do Mobile A
WITH target_user AS (
    SELECT user_id, latitude, longitude, location
    FROM user_locations
    WHERE user_id = '8bc22b4a-ad8a-4365-950c-bfd5fc7ec744'  -- Mobile A
)
SELECT 
    ul.user_id,
    ul.device_id,
    ul.latitude,
    ul.longitude,
    ul.is_anonymous,
    ST_Distance(ul.location, tu.location) as distancia_metros,
    ul.last_update,
    NOW() - ul.last_update as idade
FROM user_locations ul, target_user tu
WHERE ST_DWithin(
    ul.location,
    tu.location,
    5000  -- 5km radius
)
AND ul.last_update > NOW() - INTERVAL '30 seconds'
AND ul.user_id != tu.user_id  -- Exclui o próprio usuário
ORDER BY distancia_metros
LIMIT 101;

-- ============================================================================
-- 5. SIMULAR QUERY DE NEARBY USERS PARA USER B
-- ============================================================================
-- Substitua o UUID pelo device_id do Mobile B
WITH target_user AS (
    SELECT user_id, latitude, longitude, location
    FROM user_locations
    WHERE user_id = 'c952b59e-dc44-4ec5-a944-2d8323b6ba5a'  -- Mobile B
)
SELECT 
    ul.user_id,
    ul.device_id,
    ul.latitude,
    ul.longitude,
    ul.is_anonymous,
    ST_Distance(ul.location, tu.location) as distancia_metros,
    ul.last_update,
    NOW() - ul.last_update as idade
FROM user_locations ul, target_user tu
WHERE ST_DWithin(
    ul.location,
    tu.location,
    5000  -- 5km radius
)
AND ul.last_update > NOW() - INTERVAL '30 seconds'
AND ul.user_id != tu.user_id  -- Exclui o próprio usuário
ORDER BY distancia_metros
LIMIT 101;

-- ============================================================================
-- 6. VERIFICAR SE USER_ID É UUID VÁLIDO
-- ============================================================================
SELECT 
    user_id,
    device_id,
    is_anonymous,
    CASE 
        WHEN user_id::text ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$' 
        THEN '✅ UUID VÁLIDO'
        ELSE '❌ UUID INVÁLIDO'
    END as formato_uuid,
    last_update
FROM user_locations
WHERE is_anonymous = true
  AND last_update > NOW() - INTERVAL '5 minutes'
ORDER BY last_update DESC;

-- ============================================================================
-- 7. HISTÓRICO DE LOCALIZAÇÕES (MIGRADO PARA REDIS)
-- ============================================================================
-- Location history was migrated to Redis with automatic TTL.
-- Use Redis CLI to query history:
--   ZRANGEBYSCORE location:history:<user_id> <start_timestamp> <end_timestamp> WITHSCORES
--
-- Example: Get last 24 hours for user 8bc22b4a-ad8a-4365-950c-bfd5fc7ec744
--   ZRANGEBYSCORE location:history:8bc22b4a-ad8a-4365-950c-bfd5fc7ec744 $(date -u -d '24 hours ago' +%s) $(date -u +%s) WITHSCORES
--
-- See docs/REDIS_LOCATION_HISTORY.md for more details.

-- ============================================================================
-- 8. VERIFICAR TODAS AS LOCALIZAÇÕES ATIVAS (MAPA GERAL)
-- ============================================================================
SELECT 
    user_id,
    device_id,
    latitude,
    longitude,
    is_anonymous,
    avatar_id,
    color,
    speed,
    heading,
    last_update,
    NOW() - last_update as idade,
    CASE 
        WHEN last_update > NOW() - INTERVAL '30 seconds' THEN '✅ ATIVO'
        ELSE '❌ STALE'
    END as status
FROM user_locations
ORDER BY last_update DESC;

-- ============================================================================
-- 9. VERIFICAR ÍNDICES (PERFORMANCE)
-- ============================================================================
SELECT
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'user_locations'
ORDER BY indexname;

-- ============================================================================
-- 10. ESTATÍSTICAS DE USO
-- ============================================================================
SELECT 
    COUNT(*) as total_usuarios,
    COUNT(*) FILTER (WHERE is_anonymous = true) as total_anonimos,
    COUNT(*) FILTER (WHERE is_anonymous = false) as total_autenticados,
    COUNT(*) FILTER (WHERE last_update > NOW() - INTERVAL '30 seconds') as ativos_30s,
    COUNT(*) FILTER (WHERE is_anonymous = true AND last_update > NOW() - INTERVAL '30 seconds') as anonimos_ativos,
    COUNT(*) FILTER (WHERE is_anonymous = false AND last_update > NOW() - INTERVAL '30 seconds') as autenticados_ativos
FROM user_locations;
