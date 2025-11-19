CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS postgis;

DO
$$
    BEGIN
        CREATE TYPE entity_type AS ENUM ('erce', 'erfce');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;


DO
$$
    BEGIN
        CREATE TYPE report_status AS ENUM ('pending', 'verified', 'resolved', 'rejected');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

DO
$$
    BEGIN
        CREATE TYPE role_type AS ENUM ('citizen', 'erce', 'erfce', 'admin');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

DO
$$
    BEGIN
        CREATE TYPE alert_severity AS ENUM ('low', 'medium', 'high', 'critical');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

DO
$$
    BEGIN
        CREATE TYPE alert_status AS ENUM ('active', 'resolved', 'expired');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

CREATE TYPE notification_type AS ENUM ('alert', 'report');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    phone TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    alert_radius_meters INT DEFAULT 1000,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_code TEXT,
    email_verification_expires_at TIMESTAMP,
    nif TEXT UNIQUE,
    province TEXT,
    municipality TEXT,
    neighborhood TEXT,
    address TEXT,
    zip_code TEXT,
    country TEXT,
    last_login TIMESTAMP,
    home_address_name VARCHAR(255),
    home_address_address TEXT,
    home_address_lat DOUBLE PRECISION,
    home_address_lon DOUBLE PRECISION,
    work_address_name VARCHAR(255),
    work_address_address TEXT,
    work_address_lat DOUBLE PRECISION,
    work_address_lon DOUBLE PRECISION,
    failed_attempts INT DEFAULT 0,
    locked_until TIMESTAMP,
    device_fcm_token TEXT,
    device_language TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE roles
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name        TEXT UNIQUE NOT NULL,
    priority    INT         NOT NULL DEFAULT 0,
    description TEXT
);

CREATE TABLE IF NOT EXISTS permissions
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource TEXT NOT NULL,
    action   TEXT NOT NULL,
    code     TEXT GENERATED ALWAYS AS (LOWER(resource || ':' || action)) STORED,
    UNIQUE (resource, action)
);

CREATE TABLE IF NOT EXISTS role_permissions
(
    role_id       UUID      NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id UUID      NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    granted_at    TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_roles
(
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id     uuid NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    assigned_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, role_id)
);

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

CREATE TABLE risk_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    default_radius_meters INT DEFAULT 500,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE risk_topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    risk_type_id UUID NOT NULL REFERENCES risk_types(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (risk_type_id, name)
);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    risk_type_id UUID NOT NULL REFERENCES risk_types(id),
    risk_topic_id UUID REFERENCES risk_topics(id),
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    province TEXT,
    municipality TEXT,
    neighborhood TEXT,
    address TEXT,
    image_url TEXT,
    status report_status DEFAULT 'pending',
    reviewed_by UUID REFERENCES users(id),  -- ERFCE or ERCE user who reviewed
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    risk_type_id UUID NOT NULL REFERENCES risk_types(id),
    risk_topic_id UUID REFERENCES risk_topics(id),
    message TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    province VARCHAR(100) NULL,
    municipality VARCHAR(100) NULL,
    neighborhood VARCHAR(100) NULL,
    address TEXT,
    radius_meters INT NOT NULL DEFAULT 500,
    severity alert_severity DEFAULT 'medium',
    status alert_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    resolved_at TIMESTAMP
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type notification_type NOT NULL,
    reference_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sent_at TIMESTAMP DEFAULT NOW(),
    seen_at TIMESTAMP,
    UNIQUE(type, reference_id, user_id)
);

CREATE TABLE entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    entity_type entity_type NOT NULL,
    province TEXT,
    municipality TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    contact_email TEXT,
    contact_phone TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Users Table and Indexes
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


CREATE INDEX IF NOT EXISTS idx_users_email_trgm ON users USING gin (email gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_geo ON users (latitude, longitude);

-- Roles and Permissions Table and Indexes
CREATE INDEX IF NOT EXISTS idx_roles_priority ON roles (priority DESC);
CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions (role_id);

-- Reports Table and Indexes
CREATE TRIGGER set_timestamp_reports
    BEFORE UPDATE ON reports
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_reports_geo ON reports USING GIST (geography(ST_MakePoint(longitude, latitude)));

CREATE INDEX IF NOT EXISTS idx_reports_user ON reports(user_id);

-- Anonymous Sessions Table and Indexes
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


-- Users Table - Saved Locations Indexes
CREATE INDEX IF NOT EXISTS idx_users_home_location ON users USING GIST (
    ST_SetSRID(ST_MakePoint(home_address_lon, home_address_lat), 4326)
) WHERE home_address_lat IS NOT NULL AND home_address_lon IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_work_location ON users USING GIST (
    ST_SetSRID(ST_MakePoint(work_address_lon, work_address_lat), 4326)
) WHERE work_address_lat IS NOT NULL AND work_address_lon IS NOT NULL;

-- Location Sharings Table and Indexes
CREATE INDEX idx_location_sharings_token ON location_sharings(token);
CREATE INDEX idx_location_sharings_user_id ON location_sharings(user_id);
CREATE INDEX idx_location_sharings_anonymous_session_id ON location_sharings(anonymous_session_id);
CREATE INDEX idx_location_sharings_device_id ON location_sharings(device_id);
CREATE INDEX idx_location_sharings_expires_at ON location_sharings(expires_at);
CREATE INDEX idx_location_sharings_is_active ON location_sharings(is_active);
