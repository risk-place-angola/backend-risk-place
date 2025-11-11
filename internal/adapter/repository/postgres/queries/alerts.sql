-- name: CreateAlert :exec
INSERT INTO alerts (
    id, created_by, risk_type_id, risk_topic_id, message,
    latitude, longitude, province, municipality, neighborhood, address,
    radius_meters, severity, expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: ListActiveAlerts :many
SELECT * FROM alerts
WHERE status = 'active' AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY created_at DESC;

-- name: ResolveAlert :exec
UPDATE alerts
SET status = 'resolved', resolved_at = NOW()
WHERE id = $1;

-- name: ExpireAlert :exec
UPDATE alerts
SET status = 'expired'
WHERE id = $1;
