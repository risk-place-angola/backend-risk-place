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
SELECT 
    r.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM reports r
LEFT JOIN risk_types rt ON r.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
WHERE r.status = $1 AND r.is_private = FALSE AND rt.is_enabled = TRUE
ORDER BY r.created_at DESC;

-- name: ListReportsByUser :many
SELECT 
    r.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM reports r
LEFT JOIN risk_types rt ON r.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
WHERE r.user_id = $1 
ORDER BY r.created_at DESC;

-- name: GetReportByID :one
SELECT 
    r.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM reports r
LEFT JOIN risk_types rt ON r.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
WHERE r.id = $1 AND rt.is_enabled = TRUE;

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
    r.id, r.user_id, r.risk_type_id, r.risk_topic_id, r.description,
    r.latitude, r.longitude, r.province, r.municipality, r.neighborhood,
    r.address, r.image_url, r.status, r.reviewed_by, r.resolved_at,
    r.verification_count, r.rejection_count, r.expires_at, r.is_private,
    r.created_at, r.updated_at,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM reports r
LEFT JOIN risk_types rt ON r.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
WHERE r.id = ANY($1::uuid[]) AND r.is_private = FALSE AND rt.is_enabled = TRUE
ORDER BY r.created_at DESC;

-- name: CreateReportNotification :exec
INSERT INTO notifications (type, reference_id, user_id) VALUES ($1, $2, $3)
ON CONFLICT (type, reference_id, user_id) DO NOTHING;

-- name: ListReportsWithPagination :many
SELECT
    r.id, r.user_id, r.risk_type_id, r.risk_topic_id, r.description,
    r.latitude, r.longitude, r.province, r.municipality, r.neighborhood,
    r.address, r.image_url, r.status, r.reviewed_by, r.resolved_at,
    r.verification_count, r.rejection_count, r.expires_at, r.is_private,
    r.created_at, r.updated_at,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM reports r
LEFT JOIN risk_types rt ON r.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
WHERE (sqlc.narg('status')::text IS NULL OR r.status = sqlc.narg('status')::report_status) AND r.is_private = FALSE AND rt.is_enabled = TRUE
ORDER BY
    CASE WHEN $1 = 'desc' THEN r.created_at END DESC,
    CASE WHEN $1 = 'asc' THEN r.created_at END ASC
LIMIT $2 OFFSET $3;

-- name: CountReports :one
SELECT COUNT(*) FROM reports
WHERE (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::report_status) AND is_private = FALSE;

-- name: AddUserReportVote :exec
INSERT INTO report_votes (report_id, user_id, vote_type)
VALUES ($1, $2, $3)
ON CONFLICT (report_id, user_id) DO UPDATE SET vote_type = EXCLUDED.vote_type, created_at = NOW();

-- name: AddAnonymousReportVote :exec
INSERT INTO report_votes (report_id, anonymous_session_id, vote_type)
VALUES ($1, $2, $3)
ON CONFLICT (report_id, anonymous_session_id) DO UPDATE SET vote_type = EXCLUDED.vote_type, created_at = NOW();

-- name: RemoveUserVote :exec
DELETE FROM report_votes WHERE report_id = $1 AND user_id = $2;

-- name: RemoveAnonymousVote :exec
DELETE FROM report_votes WHERE report_id = $1 AND anonymous_session_id = $2;

-- name: GetUserVote :one
SELECT * FROM report_votes WHERE report_id = $1 AND user_id = $2;

-- name: GetAnonymousVote :one
SELECT * FROM report_votes WHERE report_id = $1 AND anonymous_session_id = $2;

-- name: UpdateVerificationCounts :exec
UPDATE reports 
SET verification_count = $2, rejection_count = $3, updated_at = NOW()
WHERE id = $1;

-- name: FindDuplicateReports :many
SELECT 
    r.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM reports r
LEFT JOIN risk_types rt ON r.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
WHERE r.risk_type_id = $1
  AND r.status = 'pending'
  AND r.is_private = FALSE
  AND rt.is_enabled = TRUE
  AND r.created_at > $2
  AND ll_to_earth(r.latitude, r.longitude) <@
      earth_box(ll_to_earth($3, $4), $5)
ORDER BY r.created_at DESC;

-- name: ExpireOldReports :exec
UPDATE reports
SET status = 'rejected', updated_at = NOW()
WHERE status = 'pending'
  AND expires_at IS NOT NULL
  AND expires_at < $1;

-- name: UpdateUserTrustScore :exec
UPDATE users SET trust_score = $2, updated_at = NOW() WHERE id = $1;

-- name: IncrementUserReportsSubmitted :exec
UPDATE users SET reports_submitted = reports_submitted + 1, updated_at = NOW() WHERE id = $1;

-- name: IncrementUserReportsVerified :exec
UPDATE users SET reports_verified = reports_verified + 1, updated_at = NOW() WHERE id = $1;