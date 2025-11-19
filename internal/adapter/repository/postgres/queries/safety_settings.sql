-- name: GetSafetySettingsByUserID :one
SELECT * FROM user_safety_settings WHERE user_id = $1 LIMIT 1;

-- name: UpsertSafetySettings :exec
INSERT INTO user_safety_settings (
    id, user_id,
    notifications_enabled, notification_alert_types, notification_alert_radius_mins,
    notification_report_types, notification_report_radius_mins,
    location_sharing_enabled, location_history_enabled,
    profile_visibility, anonymous_reports, show_online_status,
    auto_alerts_enabled, danger_zones_enabled, time_based_alerts_enabled,
    high_risk_start_time, high_risk_end_time,
    night_mode_enabled, night_mode_start_time, night_mode_end_time,
    created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
ON CONFLICT (user_id)
DO UPDATE SET
    notifications_enabled = EXCLUDED.notifications_enabled,
    notification_alert_types = EXCLUDED.notification_alert_types,
    notification_alert_radius_mins = EXCLUDED.notification_alert_radius_mins,
    notification_report_types = EXCLUDED.notification_report_types,
    notification_report_radius_mins = EXCLUDED.notification_report_radius_mins,
    location_sharing_enabled = EXCLUDED.location_sharing_enabled,
    location_history_enabled = EXCLUDED.location_history_enabled,
    profile_visibility = EXCLUDED.profile_visibility,
    anonymous_reports = EXCLUDED.anonymous_reports,
    show_online_status = EXCLUDED.show_online_status,
    auto_alerts_enabled = EXCLUDED.auto_alerts_enabled,
    danger_zones_enabled = EXCLUDED.danger_zones_enabled,
    time_based_alerts_enabled = EXCLUDED.time_based_alerts_enabled,
    high_risk_start_time = EXCLUDED.high_risk_start_time,
    high_risk_end_time = EXCLUDED.high_risk_end_time,
    night_mode_enabled = EXCLUDED.night_mode_enabled,
    night_mode_start_time = EXCLUDED.night_mode_start_time,
    night_mode_end_time = EXCLUDED.night_mode_end_time,
    updated_at = EXCLUDED.updated_at;

-- Anonymous User Queries

-- name: GetSafetySettingsByAnonymousSessionID :one
SELECT * FROM user_safety_settings 
WHERE anonymous_session_id = $1 AND device_id = $2 
LIMIT 1;

-- name: UpsertAnonymousSafetySettings :exec
INSERT INTO user_safety_settings (
    id, anonymous_session_id, device_id,
    notifications_enabled, notification_alert_types, notification_alert_radius_mins,
    notification_report_types, notification_report_radius_mins,
    location_sharing_enabled, location_history_enabled,
    profile_visibility, anonymous_reports, show_online_status,
    auto_alerts_enabled, danger_zones_enabled, time_based_alerts_enabled,
    high_risk_start_time, high_risk_end_time,
    night_mode_enabled, night_mode_start_time, night_mode_end_time,
    created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
ON CONFLICT (device_id)
DO UPDATE SET
    notifications_enabled = EXCLUDED.notifications_enabled,
    notification_alert_types = EXCLUDED.notification_alert_types,
    notification_alert_radius_mins = EXCLUDED.notification_alert_radius_mins,
    notification_report_types = EXCLUDED.notification_report_types,
    notification_report_radius_mins = EXCLUDED.notification_report_radius_mins,
    location_sharing_enabled = EXCLUDED.location_sharing_enabled,
    location_history_enabled = EXCLUDED.location_history_enabled,
    profile_visibility = EXCLUDED.profile_visibility,
    anonymous_reports = EXCLUDED.anonymous_reports,
    show_online_status = EXCLUDED.show_online_status,
    auto_alerts_enabled = EXCLUDED.auto_alerts_enabled,
    danger_zones_enabled = EXCLUDED.danger_zones_enabled,
    time_based_alerts_enabled = EXCLUDED.time_based_alerts_enabled,
    high_risk_start_time = EXCLUDED.high_risk_start_time,
    high_risk_end_time = EXCLUDED.high_risk_end_time,
    night_mode_enabled = EXCLUDED.night_mode_enabled,
    night_mode_start_time = EXCLUDED.night_mode_start_time,
    night_mode_end_time = EXCLUDED.night_mode_end_time,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteSafetySettingsByAnonymousSessionID :exec
DELETE FROM user_safety_settings 
WHERE anonymous_session_id = $1 AND device_id = $2;
