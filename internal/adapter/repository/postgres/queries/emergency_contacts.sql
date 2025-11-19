-- name: CreateEmergencyContact :exec
INSERT INTO emergency_contacts (id, user_id, name, phone, relation, is_priority, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetEmergencyContactByID :one
SELECT * FROM emergency_contacts WHERE id = $1 LIMIT 1;

-- name: GetEmergencyContactsByUserID :many
SELECT * FROM emergency_contacts WHERE user_id = $1 ORDER BY is_priority DESC, created_at ASC;

-- name: GetEmergencyContactByUserIDAndID :one
SELECT * FROM emergency_contacts WHERE user_id = $1 AND id = $2 LIMIT 1;

-- name: GetPriorityEmergencyContactsByUserID :many
SELECT * FROM emergency_contacts WHERE user_id = $1 AND is_priority = true ORDER BY created_at ASC;

-- name: CountPriorityEmergencyContactsByUserID :one
SELECT COUNT(*) FROM emergency_contacts WHERE user_id = $1 AND is_priority = true;

-- name: UpdateEmergencyContact :exec
UPDATE emergency_contacts
SET name = $2, phone = $3, relation = $4, is_priority = $5, updated_at = $6
WHERE id = $1;

-- name: DeleteEmergencyContact :exec
DELETE FROM emergency_contacts WHERE id = $1;

-- name: DeleteEmergencyContactByUserIDAndID :exec
DELETE FROM emergency_contacts WHERE user_id = $1 AND id = $2;
