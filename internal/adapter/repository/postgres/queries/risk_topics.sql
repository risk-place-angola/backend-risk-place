-- name: CreateRiskTopic :one
INSERT INTO risk_topics (risk_type_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRiskTopics :many
SELECT * FROM risk_topics ORDER BY created_at DESC;

-- name: ListRiskTopicsByType :many
SELECT * FROM risk_topics WHERE risk_type_id = $1 ORDER BY created_at DESC;

-- name: GetRiskTopicByID :one
SELECT * FROM risk_topics WHERE id = $1;

-- name: DeleteRiskTopic :exec
DELETE FROM risk_topics WHERE id = $1;