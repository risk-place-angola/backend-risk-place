-- name: CreateRiskType :exec
INSERT INTO risk_types (name, description, default_radius_meters)
VALUES ($1, $2, $3);

-- name: ListRiskTypes :many
SELECT * FROM risk_types WHERE is_enabled = TRUE ORDER BY created_at DESC;

-- name: GetRiskTypeByID :one
SELECT * FROM risk_types WHERE id = $1;

-- name: UpdateRiskType :exec
UPDATE risk_types
SET name = $2, description = $3, default_radius_meters = $4, updated_at = NOW()
WHERE id = $1;

-- name: UpdateRiskTypeIcon :exec
UPDATE risk_types
SET icon_path = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteRiskType :exec
DELETE FROM risk_types WHERE id = $1;

-- name: UpdateRiskTypeIsEnabled :exec
UPDATE risk_types
SET is_enabled = $2, updated_at = NOW()
WHERE id = $1;