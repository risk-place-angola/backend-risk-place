CREATE TABLE IF NOT EXISTS location_sharings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    anonymous_session_id UUID REFERENCES anonymous_sessions(id) ON DELETE CASCADE,
    device_id TEXT,
    owner_name VARCHAR(255),
    token VARCHAR(255) NOT NULL UNIQUE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    duration_minutes INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    last_updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT location_sharing_owner_check CHECK (
        (user_id IS NOT NULL AND anonymous_session_id IS NULL AND device_id IS NULL) OR
        (user_id IS NULL AND anonymous_session_id IS NOT NULL AND device_id IS NOT NULL)
    )
);

CREATE INDEX idx_location_sharings_token ON location_sharings(token);
CREATE INDEX idx_location_sharings_user_id ON location_sharings(user_id);
CREATE INDEX idx_location_sharings_anonymous_session_id ON location_sharings(anonymous_session_id);
CREATE INDEX idx_location_sharings_device_id ON location_sharings(device_id);
CREATE INDEX idx_location_sharings_expires_at ON location_sharings(expires_at);
CREATE INDEX idx_location_sharings_is_active ON location_sharings(is_active);
