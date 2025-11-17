CREATE TABLE anonymous_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id TEXT UNIQUE NOT NULL,
    device_fcm_token TEXT,
    device_platform TEXT,
    device_model TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    alert_radius_meters INT DEFAULT 1000,
    device_language TEXT DEFAULT 'pt',
    last_seen TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_anonymous_sessions_location ON anonymous_sessions(latitude, longitude);
CREATE INDEX idx_anonymous_sessions_device_id ON anonymous_sessions(device_id);
CREATE INDEX idx_anonymous_sessions_last_seen ON anonymous_sessions(last_seen);

CREATE TRIGGER set_timestamp_anonymous_sessions
    BEFORE UPDATE ON anonymous_sessions
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE OR REPLACE FUNCTION cleanup_old_anonymous_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM anonymous_sessions
    WHERE last_seen < NOW() - INTERVAL '30 days';
END;
$$ LANGUAGE plpgsql;
