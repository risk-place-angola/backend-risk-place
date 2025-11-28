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
    earth_distance(
        ll_to_earth(latitude, longitude),
        ll_to_earth(-8.8383, 13.2344)
    )::int as distance_meters,
    last_update,
    EXTRACT(EPOCH FROM (NOW() - last_update)) as seconds_ago
FROM user_locations
WHERE ll_to_earth(latitude, longitude) <@
      earth_box(ll_to_earth(-8.8383, 13.2344), 5000)
AND last_update > NOW() - INTERVAL '30 seconds'
ORDER BY earth_distance(
    ll_to_earth(latitude, longitude),
    ll_to_earth(-8.8383, 13.2344)
) ASC;
