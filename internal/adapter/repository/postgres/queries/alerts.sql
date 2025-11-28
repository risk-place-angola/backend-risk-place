-- name: CreateAlert :exec
INSERT INTO alerts (
    id, created_by, risk_type_id, risk_topic_id, message,
    latitude, longitude, province, municipality, neighborhood, address,
    radius_meters, severity, expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: GetAlertByID :one
SELECT 
    a.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM alerts a
LEFT JOIN risk_types rt ON a.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON a.risk_topic_id = rtopic.id
WHERE a.id = $1 AND rt.is_enabled = TRUE
LIMIT 1;

-- name: GetAlertsByUserID :many
SELECT 
    a.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM alerts a
LEFT JOIN risk_types rt ON a.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON a.risk_topic_id = rtopic.id
WHERE a.created_by = $1 AND rt.is_enabled = TRUE
ORDER BY a.created_at DESC;

-- name: GetSubscribedAlerts :many
SELECT 
    a.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM alerts a
INNER JOIN alert_subscriptions s ON a.id = s.alert_id
LEFT JOIN risk_types rt ON a.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON a.risk_topic_id = rtopic.id
WHERE s.user_id = $1 AND rt.is_enabled = TRUE
ORDER BY s.subscribed_at DESC;

-- name: ListActiveAlerts :many
SELECT 
    a.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM alerts a
LEFT JOIN risk_types rt ON a.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON a.risk_topic_id = rtopic.id
WHERE a.status = 'active' AND (a.expires_at IS NULL OR a.expires_at > NOW()) AND rt.is_enabled = TRUE
ORDER BY a.created_at DESC;

-- name: UpdateAlert :exec
UPDATE alerts
SET message = $2, severity = $3, radius_meters = $4
WHERE id = $1 AND created_by = $5;

-- name: DeleteAlert :exec
DELETE FROM alerts WHERE id = $1 AND created_by = $2;

-- name: ResolveAlert :exec
UPDATE alerts
SET status = 'resolved', resolved_at = NOW()
WHERE id = $1;

-- name: ExpireAlert :exec
UPDATE alerts
SET status = 'expired'
WHERE id = $1;

-- name: SubscribeToAlert :exec
INSERT INTO alert_subscriptions (id, alert_id, user_id, subscribed_at)
VALUES ($1, $2, $3, $4);

-- name: IsUserSubscribedToAlert :one
SELECT EXISTS (
    SELECT 1 FROM alert_subscriptions 
    WHERE alert_id = $1 AND user_id = $2
) AS subscribed;

-- name: UnsubscribeFromAlert :exec
DELETE FROM alert_subscriptions WHERE alert_id = $1 AND user_id = $2;

-- name: IsUserSubscribed :one
SELECT EXISTS(
    SELECT 1 FROM alert_subscriptions
    WHERE alert_id = $1 AND user_id = $2
) AS is_subscribed;

-- name: CountAlertSubscribers :one
SELECT COUNT(*) FROM alert_subscriptions WHERE alert_id = $1;

-- Anonymous User Queries

-- name: CreateAnonymousAlert :exec
INSERT INTO alerts (
    id, anonymous_session_id, device_id, risk_type_id, risk_topic_id, message,
    latitude, longitude, province, municipality, neighborhood, address,
    radius_meters, severity, expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

-- name: GetAlertsByAnonymousSessionID :many
SELECT 
    a.*,
    rt.name as risk_type_name,
    rt.icon_path as risk_type_icon_path,
    rtopic.name as risk_topic_name,
    rtopic.icon_path as risk_topic_icon_path
FROM alerts a
LEFT JOIN risk_types rt ON a.risk_type_id = rt.id
LEFT JOIN risk_topics rtopic ON a.risk_topic_id = rtopic.id
WHERE a.anonymous_session_id = $1 AND a.device_id = $2 AND rt.is_enabled = TRUE
ORDER BY a.created_at DESC;

-- name: UpdateAlertAnonymousToUser :exec
UPDATE alerts
SET created_by = $2, anonymous_session_id = NULL, device_id = NULL
WHERE anonymous_session_id = $1;

-- name: SubscribeAnonymousToAlert :exec
INSERT INTO alert_subscriptions (id, alert_id, anonymous_session_id, device_id, subscribed_at)
VALUES ($1, $2, $3, $4, $5);

-- name: IsAnonymousSubscribedToAlert :one
SELECT EXISTS (
    SELECT 1 FROM alert_subscriptions 
    WHERE alert_id = $1 AND device_id = $2
) AS subscribed;

-- name: UnsubscribeAnonymousFromAlert :exec
DELETE FROM alert_subscriptions WHERE alert_id = $1 AND anonymous_session_id = $2 AND device_id = $3;

-- name: GetSubscribedAlertsAnonymous :many
SELECT a.* FROM alerts a
INNER JOIN alert_subscriptions s ON a.id = s.alert_id
LEFT JOIN risk_types rt ON a.risk_type_id = rt.id
WHERE s.anonymous_session_id = $1 AND s.device_id = $2 AND rt.is_enabled = TRUE
ORDER BY s.subscribed_at DESC;

-- name: IsAnonymousSubscribed :one
SELECT EXISTS(
    SELECT 1 FROM alert_subscriptions
    WHERE alert_id = $1 AND anonymous_session_id = $2 AND device_id = $3
) AS is_subscribed;

-- name: UpdateAlertSubscriptionsAnonymousToUser :exec
UPDATE alert_subscriptions
SET user_id = $2, anonymous_session_id = NULL, device_id = NULL
WHERE anonymous_session_id = $1;
