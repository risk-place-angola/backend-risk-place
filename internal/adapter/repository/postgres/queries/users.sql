-- name: CreateUser :exec
INSERT INTO users (id, name, email, password, phone)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

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