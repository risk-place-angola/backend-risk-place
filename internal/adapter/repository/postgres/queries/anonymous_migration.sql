-- ============================================================================
-- SQLC QUERIES FOR ANONYMOUS TO AUTHENTICATED USER MIGRATION
-- ============================================================================
-- File: internal/adapter/repository/postgres/queries/anonymous_migration.sql
-- Purpose: Type-safe queries for migrating anonymous user data
-- ============================================================================

-- name: MigrateAlertsToUser :execrows
-- Migra todos os alertas de uma sessão anônima para um usuário autenticado
UPDATE alerts
SET created_by = sqlc.arg(user_id)::UUID,
    anonymous_session_id = NULL,
    device_id = NULL
WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID
  AND created_by IS NULL;

-- name: MigrateSubscriptionsToUser :execrows
-- Migra todas as subscrições de uma sessão anônima para um usuário autenticado
UPDATE alert_subscriptions
SET user_id = sqlc.arg(user_id)::UUID,
    anonymous_session_id = NULL,
    device_id = NULL
WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID
  AND user_id IS NULL;

-- name: MigrateSafetySettingsToUser :execrows
-- Migra configurações de segurança de uma sessão anônima para um usuário autenticado
-- NOTA: Apenas migra se o usuário NÃO tiver configurações
UPDATE user_safety_settings
SET user_id = sqlc.arg(user_id)::UUID,
    anonymous_session_id = NULL,
    device_id = NULL
WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID
  AND user_id IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM user_safety_settings
    WHERE user_id = sqlc.arg(user_id)::UUID
  );

-- name: MigrateLocationSharingsToUser :execrows
-- Migra compartilhamentos de localização de uma sessão anônima para um usuário autenticado
UPDATE location_sharings
SET user_id = sqlc.arg(user_id)::UUID,
    anonymous_session_id = NULL,
    device_id = NULL
WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID
  AND user_id IS NULL;

-- name: MarkAnonymousSessionAsMigrated :exec
-- Marca uma sessão anônima como migrada
UPDATE anonymous_sessions
SET migrated_to_user_id = sqlc.arg(user_id)::UUID,
    migrated_at = NOW(),
    is_active = false,
    updated_at = NOW()
WHERE id = sqlc.arg(anonymous_session_id)::UUID;

-- name: GetAnonymousDataCounts :one
-- Retorna contadores de dados anônimos antes da migração
SELECT
    COALESCE((SELECT COUNT(*) FROM alerts WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID), 0) AS alerts_count,
    COALESCE((SELECT COUNT(*) FROM alert_subscriptions WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID), 0) AS subscriptions_count,
    COALESCE((SELECT COUNT(*) FROM user_safety_settings WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID), 0) AS settings_count,
    COALESCE((SELECT COUNT(*) FROM location_sharings WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID), 0) AS location_sharings_count;

-- name: CreateDeviceUserMapping :exec
-- Cria um novo mapeamento device ↔ user
INSERT INTO device_user_mappings (
    id,
    device_id,
    anonymous_session_id,
    user_id,
    mapped_at,
    is_active
) VALUES (
    sqlc.arg(id)::UUID,
    sqlc.arg(device_id),
    sqlc.arg(anonymous_session_id)::UUID,
    sqlc.arg(user_id)::UUID,
    NOW(),
    true
);

-- name: GetActiveDeviceUserMapping :one
-- Retorna o mapeamento ativo de um device
SELECT * FROM device_user_mappings
WHERE device_id = sqlc.arg(device_id)
  AND is_active = true
LIMIT 1;

-- name: GetDeviceUserMappingsByUserID :many
-- Retorna todos os devices vinculados a um usuário
SELECT * FROM device_user_mappings
WHERE user_id = sqlc.arg(user_id)::UUID
ORDER BY mapped_at DESC;

-- name: DeactivateDeviceUserMapping :exec
-- Marca um mapeamento como inativo
UPDATE device_user_mappings
SET is_active = false,
    unmapped_at = NOW()
WHERE id = sqlc.arg(id)::UUID;

-- name: CreateAnonymousMigration :exec
-- Cria um novo log de migração
INSERT INTO anonymous_user_migrations (
    id,
    anonymous_session_id,
    device_id,
    user_id,
    migration_type,
    started_at
) VALUES (
    sqlc.arg(id)::UUID,
    sqlc.arg(anonymous_session_id)::UUID,
    sqlc.arg(device_id),
    sqlc.arg(user_id)::UUID,
    sqlc.arg(migration_type),
    NOW()
);

-- name: UpdateMigrationCounters :exec
-- Atualiza os contadores de uma migração
UPDATE anonymous_user_migrations
SET alerts_migrated = sqlc.arg(alerts_migrated),
    subscriptions_migrated = sqlc.arg(subscriptions_migrated),
    settings_migrated = sqlc.arg(settings_migrated)::BOOLEAN,
    location_sharings_migrated = sqlc.arg(location_sharings_migrated)
WHERE id = sqlc.arg(id)::UUID;

-- name: MarkMigrationCompleted :exec
-- Marca uma migração como concluída
UPDATE anonymous_user_migrations
SET completed_at = NOW()
WHERE id = sqlc.arg(id)::UUID;

-- name: MarkMigrationFailed :exec
-- Marca uma migração como falha
UPDATE anonymous_user_migrations
SET failed_at = NOW(),
    error_message = sqlc.arg(error_message)
WHERE id = sqlc.arg(id)::UUID;

-- name: GetMigrationsByDeviceID :many
-- Retorna todas as migrações de um device
SELECT * FROM anonymous_user_migrations
WHERE device_id = sqlc.arg(device_id)
ORDER BY started_at DESC;

-- name: GetMigrationsByUserID :many
-- Retorna todas as migrações de um usuário
SELECT * FROM anonymous_user_migrations
WHERE user_id = sqlc.arg(user_id)::UUID
ORDER BY started_at DESC;

-- name: GetMigrationByID :one
-- Retorna uma migração específica
SELECT * FROM anonymous_user_migrations
WHERE id = sqlc.arg(id)::UUID;

-- ============================================================================
-- QUERIES FOR MERGE STRATEGY (Settings Conflict Resolution)
-- ============================================================================

-- name: GetSafetySettingsByAnonymousSession :one
-- Retorna configurações de uma sessão anônima
SELECT * FROM user_safety_settings
WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID
LIMIT 1;

-- name: GetUserSafetySettings :one
-- Retorna configurações de um usuário autenticado
SELECT * FROM user_safety_settings
WHERE user_id = sqlc.arg(user_id)::UUID
LIMIT 1;

-- name: MergeSafetySettings :exec
-- Merge inteligente: prioriza configurações anônimas se mais recentes
UPDATE user_safety_settings auth
SET 
    notifications_enabled = COALESCE(anon.notifications_enabled, auth.notifications_enabled),
    notification_alert_types = COALESCE(anon.notification_alert_types, auth.notification_alert_types),
    notification_alert_radius_mins = COALESCE(anon.notification_alert_radius_mins, auth.notification_alert_radius_mins),
    notification_report_types = COALESCE(anon.notification_report_types, auth.notification_report_types),
    notification_report_radius_mins = COALESCE(anon.notification_report_radius_mins, auth.notification_report_radius_mins),
    location_sharing_enabled = COALESCE(anon.location_sharing_enabled, auth.location_sharing_enabled),
    location_history_enabled = COALESCE(anon.location_history_enabled, auth.location_history_enabled),
    profile_visibility = COALESCE(anon.profile_visibility, auth.profile_visibility),
    anonymous_reports = COALESCE(anon.anonymous_reports, auth.anonymous_reports),
    show_online_status = COALESCE(anon.show_online_status, auth.show_online_status),
    auto_alerts_enabled = COALESCE(anon.auto_alerts_enabled, auth.auto_alerts_enabled),
    danger_zones_enabled = COALESCE(anon.danger_zones_enabled, auth.danger_zones_enabled),
    time_based_alerts_enabled = COALESCE(anon.time_based_alerts_enabled, auth.time_based_alerts_enabled),
    high_risk_start_time = COALESCE(anon.high_risk_start_time, auth.high_risk_start_time),
    high_risk_end_time = COALESCE(anon.high_risk_end_time, auth.high_risk_end_time),
    night_mode_enabled = COALESCE(anon.night_mode_enabled, auth.night_mode_enabled),
    night_mode_start_time = COALESCE(anon.night_mode_start_time, auth.night_mode_start_time),
    night_mode_end_time = COALESCE(anon.night_mode_end_time, auth.night_mode_end_time),
    updated_at = NOW()
FROM user_safety_settings anon
WHERE auth.user_id = sqlc.arg(user_id)::UUID
  AND anon.anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID
  AND anon.updated_at > auth.updated_at;  -- Apenas se anon mais recente

-- name: DeleteAnonymousSafetySettings :exec
-- Remove configurações anônimas após merge
DELETE FROM user_safety_settings
WHERE anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID;

-- ============================================================================
-- QUERIES FOR ROLLBACK (Error Recovery)
-- ============================================================================

-- name: RollbackAlertsToAnonymous :execrows
-- Reverte migração de alertas (rollback)
UPDATE alerts
SET created_by = NULL,
    anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID,
    device_id = sqlc.arg(device_id)
WHERE created_by = sqlc.arg(user_id)::UUID
  AND created_at >= sqlc.arg(migration_started_at)::TIMESTAMP;

-- name: RollbackSubscriptionsToAnonymous :execrows
-- Reverte migração de subscrições (rollback)
UPDATE alert_subscriptions
SET user_id = NULL,
    anonymous_session_id = sqlc.arg(anonymous_session_id)::UUID,
    device_id = sqlc.arg(device_id)
WHERE user_id = sqlc.arg(user_id)::UUID
  AND subscribed_at >= sqlc.arg(migration_started_at)::TIMESTAMP;

-- name: RollbackAnonymousSessionMigration :exec
-- Reverte status de migração da sessão anônima
UPDATE anonymous_sessions
SET migrated_to_user_id = NULL,
    migrated_at = NULL,
    is_active = true,
    updated_at = NOW()
WHERE id = sqlc.arg(anonymous_session_id)::UUID;

-- ============================================================================
-- ANALYTICS QUERIES
-- ============================================================================

-- name: GetMigrationStats :one
-- Retorna estatísticas gerais de migrações
SELECT
    COUNT(*) AS total_migrations,
    COUNT(*) FILTER (WHERE completed_at IS NOT NULL) AS completed_count,
    COUNT(*) FILTER (WHERE failed_at IS NOT NULL) AS failed_count,
    COUNT(*) FILTER (WHERE migration_type = 'signup') AS signup_count,
    COUNT(*) FILTER (WHERE migration_type = 'login') AS login_count,
    COALESCE(AVG(alerts_migrated), 0) AS avg_alerts_per_migration,
    COALESCE(AVG(subscriptions_migrated), 0) AS avg_subscriptions_per_migration
FROM anonymous_user_migrations
WHERE started_at >= sqlc.arg(since)::TIMESTAMP;

-- name: GetRecentFailedMigrations :many
-- Retorna migrações falhadas recentes para debug
SELECT * FROM anonymous_user_migrations
WHERE failed_at IS NOT NULL
ORDER BY failed_at DESC
LIMIT $1;

-- name: GetUserMigrationHistory :many
-- Retorna histórico completo de migrações de um usuário
SELECT 
    m.*,
    s.device_platform,
    s.device_model
FROM anonymous_user_migrations m
JOIN anonymous_sessions s ON m.anonymous_session_id = s.id
WHERE m.user_id = sqlc.arg(user_id)::UUID
ORDER BY m.started_at DESC;
