-- name: CreateAlertNotification :exec
INSERT INTO notifications (type, reference_id, user_id) VALUES ($1, $2, $3)
ON CONFLICT (type, reference_id, user_id) DO NOTHING;

-- name: MarkAlertSeen :exec
UPDATE notifications
SET seen_at = NOW(),
    updated_at = NOW()
WHERE reference_id = $1 AND user_id = $2;

-- name: ListUserNotifications :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY sent_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserUnreadNotifications :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND seen_at IS NULL;
