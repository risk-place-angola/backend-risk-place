-- name: CreateEntity :one
INSERT INTO entities (
    name, entity_type, province, municipality, latitude, longitude, contact_email, contact_phone
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListEntities :many
SELECT * FROM entities ORDER BY created_at DESC;

-- name: GetEntitiesByType :many
SELECT * FROM entities WHERE entity_type = $1 ORDER BY created_at DESC;

-- name: DeleteEntity :exec
DELETE FROM entities WHERE id = $1;
