-- Drop function
DROP FUNCTION IF EXISTS cleanup_old_anonymous_sessions();

-- Drop trigger
DROP TRIGGER IF EXISTS set_timestamp_anonymous_sessions ON anonymous_sessions;

-- Drop indexes
DROP INDEX IF EXISTS idx_anonymous_sessions_last_seen;
DROP INDEX IF EXISTS idx_anonymous_sessions_device_id;
DROP INDEX IF EXISTS idx_anonymous_sessions_location;

-- Drop table
DROP TABLE IF EXISTS anonymous_sessions;
