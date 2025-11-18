-- name: CreateReport :one
INSERT INTO reports (
    user_id, risk_type_id, risk_topic_id, description,
    latitude, longitude, province, municipality,
    neighborhood, address, image_url
) VALUES (
             $1, $2, $3, $4, $5,
             $6, $7, $8, $9, $10, $11
         )
RETURNING id;

-- name: ListReportsByStatus :many
SELECT * FROM reports WHERE status = $1 ORDER BY created_at DESC;

-- name: ListReportsByUser :many
SELECT * FROM reports WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetReportByID :one
SELECT * FROM reports WHERE id = $1;

-- name: VerifyReport :exec
UPDATE reports
SET status = 'verified', reviewed_by = $2, updated_at = NOW()
WHERE id = $1;

-- name: ResolveReport :exec
UPDATE reports
SET status = 'resolved',
    resolved_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: RejectReport :exec
UPDATE reports
SET status = 'rejected',
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateReport :exec
UPDATE reports
SET description = $2,
    status = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteReport :exec
DELETE FROM reports WHERE id = $1;

-- name: UpdateReportLocation :one
UPDATE reports
SET latitude = $2,
    longitude = $3,
    address = COALESCE(NULLIF($4, ''), address),
    neighborhood = COALESCE(NULLIF($5, ''), neighborhood),
    municipality = COALESCE(NULLIF($6, ''), municipality),
    province = COALESCE(NULLIF($7, ''), province),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;

-- name: ListReportsByIDs :many
SELECT
    id, user_id, risk_type_id, risk_topic_id, description,
    latitude, longitude, province, municipality, neighborhood,
    address, image_url, status, reviewed_by, resolved_at,
    created_at, updated_at
FROM reports
WHERE id = ANY($1::uuid[])
ORDER BY created_at DESC;

-- name: CreateReportNotification :exec
INSERT INTO notifications (type, reference_id, user_id) VALUES ($1, $2, $3)
ON CONFLICT (type, reference_id, user_id) DO NOTHING;

-- name: ListReportsWithPagination :many
SELECT
    id, user_id, risk_type_id, risk_topic_id, description,
    latitude, longitude, province, municipality, neighborhood,
    address, image_url, status, reviewed_by, resolved_at,
    created_at, updated_at
FROM reports
WHERE (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::report_status)
ORDER BY
    CASE WHEN $1 = 'desc' THEN created_at END DESC,
    CASE WHEN $1 = 'asc' THEN created_at END ASC
LIMIT $2 OFFSET $3;

-- name: CountReports :one
SELECT COUNT(*) FROM reports
WHERE (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::report_status);