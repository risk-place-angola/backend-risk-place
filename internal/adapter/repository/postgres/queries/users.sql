-- name: CreateUser :exec
INSERT INTO users (id, name, email, password, phone)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByEmailOrPhone :one
SELECT * FROM users WHERE (email = $1 OR phone = $1) AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: UpdateUserLocation :exec
UPDATE users
SET latitude = $2,
    longitude = $3,
    province = $4,
    municipality = $5,
    neighborhood = $6,
    address = $7,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateAlertRadius :exec
UPDATE users SET alert_radius_meters = $2 WHERE id = $1;

-- name: MarkAccountVerified :exec
UPDATE users SET account_verified = TRUE, email_verification_code = NULL WHERE id = $1;

-- name: UpdateEmailVerificationCode :exec
UPDATE users
SET email_verification_code = $2,
    email_verification_expires_at = $3
WHERE id = $1;

-- name: ListNearbyUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
  AND (latitude IS NOT NULL AND longitude IS NOT NULL);

-- name: AddCodeToUser :exec
UPDATE users
SET email_verification_code = $2, email_verification_expires_at = $3, updated_at = $4, account_verified = false
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserDeviceInfo :exec
UPDATE users
SET device_fcm_token = $2,
    device_language = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: ListAllDeviceTokensExceptUser :many
SELECT device_fcm_token, device_language
FROM users
WHERE deleted_at IS NULL
  AND device_fcm_token IS NOT NULL
  AND id <> $1;

-- name: ListDeviceTokensByUserIDs :many
SELECT device_fcm_token, device_language
FROM users
WHERE deleted_at IS NULL
  AND device_fcm_token IS NOT NULL
  AND id IN (SELECT UNNEST($1::uuid[]));

-- name: ListDeviceTokensForAlertNotification :many
SELECT DISTINCT u.device_fcm_token, u.device_language, u.id as user_id
FROM users u
LEFT JOIN user_safety_settings s ON s.user_id = u.id
WHERE u.deleted_at IS NULL
  AND u.device_fcm_token IS NOT NULL
  AND u.id IN (SELECT UNNEST(sqlc.arg(user_ids)::uuid[]))
  AND (s.id IS NULL OR s.notifications_enabled = true)
  AND (
    s.id IS NULL OR
    sqlc.arg(severity_level)::TEXT = ANY(s.notification_alert_types) OR
    'all' = ANY(s.notification_alert_types)
  )
  AND (
    s.id IS NULL OR
    s.notification_alert_radius_mins >= sqlc.arg(distance_meters)::INT
  );

-- name: ListDeviceTokensForReportNotification :many
SELECT DISTINCT u.device_fcm_token, u.device_language, u.id as user_id
FROM users u
LEFT JOIN user_safety_settings s ON s.user_id = u.id
WHERE u.deleted_at IS NULL
  AND u.device_fcm_token IS NOT NULL
  AND u.id IN (SELECT UNNEST(sqlc.arg(user_ids)::uuid[]))
  AND (s.id IS NULL OR s.notifications_enabled = true)
  AND (
    s.id IS NULL OR
    'all' = ANY(s.notification_report_types) OR
    (sqlc.arg(is_verified)::BOOLEAN = true AND 'verified' = ANY(s.notification_report_types))
  )
  AND (
    s.id IS NULL OR
    s.notification_report_radius_mins >= sqlc.arg(distance_meters)::INT
  );

-- name: ListAnonymousTokensForAlertNotification :many
SELECT DISTINCT a.device_fcm_token, a.device_id
FROM anonymous_sessions a
LEFT JOIN user_safety_settings s ON s.device_id = a.device_id
WHERE a.device_fcm_token IS NOT NULL
  AND a.latitude IS NOT NULL
  AND a.longitude IS NOT NULL
  AND (s.id IS NULL OR s.notifications_enabled = true)
  AND (
    6371000 * acos(
      cos(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * cos(radians(a.latitude)) *
      cos(radians(a.longitude) - radians(sqlc.arg(longitude)::DOUBLE PRECISION)) +
      sin(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * sin(radians(a.latitude))
    )
  ) <= a.alert_radius_meters
  AND (
    6371000 * acos(
      cos(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * cos(radians(a.latitude)) *
      cos(radians(a.longitude) - radians(sqlc.arg(longitude)::DOUBLE PRECISION)) +
      sin(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * sin(radians(a.latitude))
    )
  ) <= sqlc.arg(radius_meters)::DOUBLE PRECISION
  AND (
    s.id IS NULL OR
    sqlc.arg(severity_level)::TEXT = ANY(s.notification_alert_types) OR
    'all' = ANY(s.notification_alert_types)
  )
  AND (
    s.id IS NULL OR
    s.notification_alert_radius_mins >= CAST((
      6371000 * acos(
        cos(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * cos(radians(a.latitude)) *
        cos(radians(a.longitude) - radians(sqlc.arg(longitude)::DOUBLE PRECISION)) +
        sin(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * sin(radians(a.latitude))
      )
    ) AS INT)
  );

-- name: ListAnonymousTokensForReportNotification :many
SELECT DISTINCT a.device_fcm_token, a.device_id
FROM anonymous_sessions a
LEFT JOIN user_safety_settings s ON s.device_id = a.device_id
WHERE a.device_fcm_token IS NOT NULL
  AND a.latitude IS NOT NULL
  AND a.longitude IS NOT NULL
  AND (s.id IS NULL OR s.notifications_enabled = true)
  AND (
    6371000 * acos(
      cos(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * cos(radians(a.latitude)) *
      cos(radians(a.longitude) - radians(sqlc.arg(longitude)::DOUBLE PRECISION)) +
      sin(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * sin(radians(a.latitude))
    )
  ) <= a.alert_radius_meters
  AND (
    6371000 * acos(
      cos(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * cos(radians(a.latitude)) *
      cos(radians(a.longitude) - radians(sqlc.arg(longitude)::DOUBLE PRECISION)) +
      sin(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * sin(radians(a.latitude))
    )
  ) <= sqlc.arg(radius_meters)::DOUBLE PRECISION
  AND (
    s.id IS NULL OR
    'all' = ANY(s.notification_report_types) OR
    (sqlc.arg(is_verified)::BOOLEAN = true AND 'verified' = ANY(s.notification_report_types))
  )
  AND (
    s.id IS NULL OR
    s.notification_report_radius_mins >= CAST((
      6371000 * acos(
        cos(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * cos(radians(a.latitude)) *
        cos(radians(a.longitude) - radians(sqlc.arg(longitude)::DOUBLE PRECISION)) +
        sin(radians(sqlc.arg(latitude)::DOUBLE PRECISION)) * sin(radians(a.latitude))
      )
    ) AS INT)
  );

-- name: UpdateUserSavedLocations :exec
UPDATE users
SET home_address_name = $2,
    home_address_address = $3,
    home_address_lat = $4,
    home_address_lon = $5,
    work_address_name = $6,
    work_address_address = $7,
    work_address_lat = $8,
    work_address_lon = $9,
    updated_at = NOW()
WHERE id = $1;