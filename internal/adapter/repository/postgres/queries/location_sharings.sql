-- name: CreateLocationSharing :exec
INSERT INTO location_sharings (
    id, user_id, anonymous_session_id, device_id, owner_name,
    token, latitude, longitude, duration_minutes, expires_at, last_updated_at, is_active
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: UpdateLocationSharing :exec
UPDATE location_sharings
SET 
    latitude = $2,
    longitude = $3,
    last_updated_at = $4,
    is_active = $5,
    updated_at = $6
WHERE id = $1;

-- name: DeleteLocationSharing :exec
DELETE FROM location_sharings
WHERE id = $1;

-- name: GetLocationSharingByID :one
SELECT * FROM location_sharings
WHERE id = $1;

-- name: GetLocationSharingByToken :one
SELECT * FROM location_sharings
WHERE token = $1;

-- name: ListActiveLocationSharingsByUserID :many
SELECT * FROM location_sharings
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: ListActiveLocationSharingsByDeviceID :many
SELECT * FROM location_sharings
WHERE device_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: ListAllLocationSharings :many
SELECT * FROM location_sharings
ORDER BY created_at DESC;

-- name: DeactivateExpiredLocationSharings :exec
UPDATE location_sharings
SET is_active = false, updated_at = NOW()
WHERE expires_at < $1 AND is_active = true;
