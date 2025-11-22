-- Script para verificar os dados de user_locations
-- Use este script para debugar problemas com nearby users

-- 1. Contar total de localizações
SELECT COUNT(*) as total_locations FROM user_locations;

-- 2. Ver todas as localizações (últimos 30 segundos)
SELECT 
    user_id,
    device_id,
    latitude,
    longitude,
    is_anonymous,
    last_update,
    EXTRACT(EPOCH FROM (NOW() - last_update)) as seconds_ago
FROM user_locations
WHERE last_update > NOW() - INTERVAL '30 seconds'
ORDER BY last_update DESC;

-- 3. Ver todas as localizações (todas)
SELECT 
    user_id,
    device_id,
    latitude,
    longitude,
    is_anonymous,
    last_update,
    EXTRACT(EPOCH FROM (NOW() - last_update)) as seconds_ago
FROM user_locations
ORDER BY last_update DESC
LIMIT 20;

-- 4. Testar query de nearby users (substitua LAT, LON pelo valor real)
-- Exemplo: Para Luanda, Angola use: latitude=-8.8383, longitude=13.2344
SELECT 
    user_id,
    latitude,
    longitude,
    is_anonymous,
    ST_Distance(
        location,
        ST_SetSRID(ST_MakePoint(13.2344, -8.8383), 4326)::geography
    ) as distance_meters,
    last_update,
    EXTRACT(EPOCH FROM (NOW() - last_update)) as seconds_ago
FROM user_locations
WHERE ST_DWithin(
    location,
    ST_SetSRID(ST_MakePoint(13.2344, -8.8383), 4326)::geography,
    5000
)
AND last_update > NOW() - INTERVAL '30 seconds'
ORDER BY location <-> ST_SetSRID(ST_MakePoint(13.2344, -8.8383), 4326)::geography;
