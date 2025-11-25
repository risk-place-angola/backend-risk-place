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
    account_verified BOOLEAN DEFAULT FALSE,
    email_verification_code TEXT,
    email_verification_expires_at TIMESTAMP,
    nif TEXT UNIQUE,
    province TEXT,
    municipality TEXT,
    neighborhood TEXT,
    address TEXT,
    zip_code TEXT,
    country TEXT,
    push_notification_enabled BOOLEAN DEFAULT true,
    sms_notification_enabled BOOLEAN DEFAULT false,
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
    trust_score INT DEFAULT 50,
    reports_submitted INT DEFAULT 0,
    reports_verified INT DEFAULT 0,
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
    push_notification_enabled BOOLEAN DEFAULT true,
    sms_notification_enabled BOOLEAN DEFAULT false,
    alert_radius_meters INT DEFAULT 1000,
    device_language TEXT DEFAULT 'pt',
    last_seen TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    migrated_to_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    migrated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE emergency_contacts (
                                    id UUID PRIMARY KEY,
                                    user_id UUID NOT NULL,
                                    name VARCHAR(255) NOT NULL,
                                    phone VARCHAR(20) NOT NULL,
                                    relation VARCHAR(50) NOT NULL CHECK (relation IN ('family', 'friend', 'colleague', 'neighbor', 'other')),
                                    is_priority BOOLEAN NOT NULL DEFAULT false,
                                    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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
    icon_path TEXT,
    default_radius_meters INT DEFAULT 500,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE risk_topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    risk_type_id UUID NOT NULL REFERENCES risk_types(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    icon_path TEXT,
    is_sensitive BOOLEAN NOT NULL DEFAULT FALSE,
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
    reviewed_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    verification_count INT DEFAULT 0,
    rejection_count INT DEFAULT 0,
    expires_at TIMESTAMP,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    anonymous_session_id UUID REFERENCES anonymous_sessions(id) ON DELETE SET NULL,
    device_id TEXT,
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
    resolved_at TIMESTAMP,
    CONSTRAINT alerts_creator_check CHECK (
        (created_by IS NOT NULL AND anonymous_session_id IS NULL AND device_id IS NULL) OR
        (created_by IS NULL AND anonymous_session_id IS NOT NULL AND device_id IS NOT NULL)
    )
);

CREATE TABLE alert_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id UUID NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    anonymous_session_id UUID REFERENCES anonymous_sessions(id) ON DELETE CASCADE,
    device_id TEXT,
    subscribed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT alert_subscriptions_subscriber_check CHECK (
        (user_id IS NOT NULL AND anonymous_session_id IS NULL AND device_id IS NULL) OR
        (user_id IS NULL AND anonymous_session_id IS NOT NULL AND device_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX idx_alert_subscriptions_unique_user ON alert_subscriptions(alert_id, user_id) WHERE user_id IS NOT NULL;
CREATE UNIQUE INDEX idx_alert_subscriptions_unique_device ON alert_subscriptions(alert_id, device_id) WHERE device_id IS NOT NULL;

CREATE TABLE report_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    anonymous_session_id UUID REFERENCES anonymous_sessions(id) ON DELETE CASCADE,
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('upvote', 'downvote')),
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT vote_user_or_anonymous CHECK (
        (user_id IS NOT NULL AND anonymous_session_id IS NULL) OR
        (user_id IS NULL AND anonymous_session_id IS NOT NULL)
    ),
    UNIQUE(report_id, user_id),
    UNIQUE(report_id, anonymous_session_id)
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
CREATE INDEX IF NOT EXISTS idx_users_trust_score ON users(trust_score);

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
CREATE INDEX IF NOT EXISTS idx_reports_status_created ON reports(status, created_at);
CREATE INDEX IF NOT EXISTS idx_reports_expires_at ON reports(expires_at);
CREATE INDEX IF NOT EXISTS idx_reports_is_private ON reports(is_private);
CREATE INDEX IF NOT EXISTS idx_reports_status_private ON reports(status, is_private);

CREATE INDEX IF NOT EXISTS idx_reports_geo ON reports USING GIST (geography(ST_MakePoint(longitude, latitude)));

CREATE INDEX IF NOT EXISTS idx_reports_user ON reports(user_id);

-- Report Votes Table and Indexes
CREATE INDEX idx_report_votes_report ON report_votes(report_id);
CREATE INDEX idx_report_votes_user ON report_votes(user_id);

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

-- Emergency Contacts Table and Indexes
CREATE INDEX idx_emergency_contacts_user_id ON emergency_contacts(user_id);
CREATE INDEX idx_emergency_contacts_priority ON emergency_contacts(user_id, is_priority) WHERE is_priority = true;

-- Alert Subscriptions Table and Indexes
CREATE INDEX idx_alert_subscriptions_alert_id ON alert_subscriptions(alert_id);
CREATE INDEX idx_alert_subscriptions_user_id ON alert_subscriptions(user_id);
CREATE INDEX idx_alert_subscriptions_subscribed_at ON alert_subscriptions(subscribed_at DESC);


-- User Safety Settings Table
CREATE TABLE user_safety_settings (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                      user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                      anonymous_session_id UUID REFERENCES anonymous_sessions(id) ON DELETE CASCADE,
                                      device_id TEXT,

    -- Notification Settings
                                      notifications_enabled BOOLEAN DEFAULT true,
                                      notification_alert_types TEXT[] DEFAULT ARRAY['high', 'critical'],
                                      notification_alert_radius_mins INT DEFAULT 1000,
                                      notification_report_types TEXT[] DEFAULT ARRAY['verified'],
                                      notification_report_radius_mins INT DEFAULT 500,

    -- Tracking Settings
                                      location_sharing_enabled BOOLEAN DEFAULT false,
                                      location_history_enabled BOOLEAN DEFAULT true,

    -- Privacy Settings
                                      profile_visibility TEXT DEFAULT 'public' CHECK (profile_visibility IN ('public', 'friends', 'private')),
                                      anonymous_reports BOOLEAN DEFAULT false,
                                      show_online_status BOOLEAN DEFAULT true,

    -- Auto Alert Settings
                                      auto_alerts_enabled BOOLEAN DEFAULT false,
                                      danger_zones_enabled BOOLEAN DEFAULT true,
                                      time_based_alerts_enabled BOOLEAN DEFAULT false,
                                      high_risk_start_time TIME DEFAULT '22:00',
                                      high_risk_end_time TIME DEFAULT '06:00',

    -- Night Mode Settings
                                      night_mode_enabled BOOLEAN DEFAULT false,
                                      night_mode_start_time TIME DEFAULT '22:00',
                                      night_mode_end_time TIME DEFAULT '06:00',

                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      
    CONSTRAINT user_safety_settings_owner_check CHECK (
        (user_id IS NOT NULL AND anonymous_session_id IS NULL AND device_id IS NULL) OR
        (user_id IS NULL AND anonymous_session_id IS NOT NULL AND device_id IS NOT NULL)
    )
);

-- User Safety Settings Indexes
CREATE UNIQUE INDEX idx_safety_settings_unique_user ON user_safety_settings(user_id) WHERE user_id IS NOT NULL;
CREATE UNIQUE INDEX idx_safety_settings_unique_device ON user_safety_settings(device_id) WHERE device_id IS NOT NULL;
CREATE INDEX idx_user_safety_settings_user_id ON user_safety_settings(user_id);
CREATE INDEX idx_safety_settings_anonymous_session_id ON user_safety_settings(anonymous_session_id) WHERE anonymous_session_id IS NOT NULL;
CREATE INDEX idx_safety_settings_device_id ON user_safety_settings(device_id) WHERE device_id IS NOT NULL;
CREATE INDEX idx_user_safety_settings_updated_at ON user_safety_settings(updated_at DESC);


CREATE INDEX IF NOT EXISTS idx_users_notification_prefs ON users(push_notification_enabled, sms_notification_enabled) WHERE push_notification_enabled = true OR sms_notification_enabled = true;
CREATE INDEX IF NOT EXISTS idx_anonymous_sessions_notification_prefs ON anonymous_sessions(push_notification_enabled, sms_notification_enabled) WHERE push_notification_enabled = true OR sms_notification_enabled = true;


-- ============================================================================
-- ANONYMOUS USER TO AUTHENTICATED USER MIGRATION SUPPORT
-- ============================================================================

-- Add indexes for alerts anonymous support
CREATE INDEX idx_alerts_anonymous_session_id ON alerts(anonymous_session_id) WHERE anonymous_session_id IS NOT NULL;
CREATE INDEX idx_alerts_device_id ON alerts(device_id) WHERE device_id IS NOT NULL;

-- Add indexes for alert_subscriptions anonymous support
CREATE INDEX idx_alert_subscriptions_anonymous_session_id ON alert_subscriptions(anonymous_session_id) WHERE anonymous_session_id IS NOT NULL;
CREATE INDEX idx_alert_subscriptions_device_id ON alert_subscriptions(device_id) WHERE device_id IS NOT NULL;

-- Add indexes for anonymous_sessions migration tracking
CREATE INDEX idx_anonymous_sessions_migrated_user ON anonymous_sessions(migrated_to_user_id) WHERE migrated_to_user_id IS NOT NULL;
CREATE INDEX idx_anonymous_sessions_is_active ON anonymous_sessions(is_active) WHERE is_active = true;

-- Add linked device to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS linked_device_id TEXT;
CREATE INDEX idx_users_linked_device_id ON users(linked_device_id) WHERE linked_device_id IS NOT NULL;

-- Device-User Mapping Table (Audit Trail and Re-linking)
CREATE TABLE device_user_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id TEXT NOT NULL,
    anonymous_session_id UUID NOT NULL REFERENCES anonymous_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mapped_at TIMESTAMP NOT NULL DEFAULT NOW(),
    unmapped_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    CONSTRAINT device_user_mapping_active_check CHECK (
        (is_active = true AND unmapped_at IS NULL) OR
        (is_active = false AND unmapped_at IS NOT NULL)
    )
);

CREATE UNIQUE INDEX idx_device_user_mappings_active ON device_user_mappings(device_id, user_id) WHERE is_active = true;
CREATE INDEX idx_device_user_mappings_device_id ON device_user_mappings(device_id);
CREATE INDEX idx_device_user_mappings_user_id ON device_user_mappings(user_id);
CREATE INDEX idx_device_user_mappings_anonymous_session_id ON device_user_mappings(anonymous_session_id);

-- Migration Audit Log
CREATE TABLE anonymous_user_migrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anonymous_session_id UUID NOT NULL REFERENCES anonymous_sessions(id) ON DELETE CASCADE,
    device_id TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Counters
    alerts_migrated INT NOT NULL DEFAULT 0,
    subscriptions_migrated INT NOT NULL DEFAULT 0,
    settings_migrated BOOLEAN NOT NULL DEFAULT false,
    location_sharings_migrated INT NOT NULL DEFAULT 0,
    
    -- Metadata
    migration_type VARCHAR(20) NOT NULL CHECK (migration_type IN ('signup', 'login', 'manual')),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    failed_at TIMESTAMP,
    error_message TEXT,
    
    CONSTRAINT migration_completion_check CHECK (
        (completed_at IS NOT NULL AND failed_at IS NULL) OR
        (completed_at IS NULL AND failed_at IS NOT NULL) OR
        (completed_at IS NULL AND failed_at IS NULL)
    )
);

CREATE INDEX idx_migrations_anonymous_session_id ON anonymous_user_migrations(anonymous_session_id);
CREATE INDEX idx_migrations_user_id ON anonymous_user_migrations(user_id);
CREATE INDEX idx_migrations_device_id ON anonymous_user_migrations(device_id);
CREATE INDEX idx_migrations_started_at ON anonymous_user_migrations(started_at DESC);

-- ============================================================================
-- FUNCTIONS AND TRIGGERS FOR ANONYMOUS USER SUPPORT
-- ============================================================================

-- Auto-expire anonymous alerts after 2 hours
CREATE OR REPLACE FUNCTION auto_expire_anonymous_alerts()
RETURNS void AS $$
BEGIN
    UPDATE alerts
    SET status = 'expired'
    WHERE anonymous_session_id IS NOT NULL
      AND status = 'active'
      AND created_at < NOW() - INTERVAL '2 hours';
END;
$$ LANGUAGE plpgsql;

-- Update anonymous_sessions.last_seen when activity occurs
CREATE OR REPLACE FUNCTION update_anonymous_session_last_seen()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE anonymous_sessions
    SET last_seen = NOW()
    WHERE id = NEW.anonymous_session_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_last_seen_on_alert
    AFTER INSERT OR UPDATE ON alerts
    FOR EACH ROW
    WHEN (NEW.anonymous_session_id IS NOT NULL)
EXECUTE FUNCTION update_anonymous_session_last_seen();

CREATE TRIGGER trigger_update_last_seen_on_subscription
    AFTER INSERT OR UPDATE ON alert_subscriptions
    FOR EACH ROW
    WHEN (NEW.anonymous_session_id IS NOT NULL)
EXECUTE FUNCTION update_anonymous_session_last_seen();

-- ============================================================================
-- USER LOCATIONS FOR NEARBY USERS FEATURE (Waze-style)
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    device_id VARCHAR(255),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    location GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography) STORED,
    speed DOUBLE PRECISION DEFAULT 0,
    heading DOUBLE PRECISION DEFAULT 0,
    avatar_id INTEGER NOT NULL,
    color VARCHAR(7) NOT NULL,
    is_anonymous BOOLEAN DEFAULT false,
    last_update TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT user_locations_user_id_key UNIQUE (user_id)
);

CREATE INDEX idx_user_locations_location ON user_locations USING GIST (location);
CREATE INDEX idx_user_locations_last_update ON user_locations (last_update);
CREATE INDEX idx_user_locations_user_id ON user_locations (user_id);
CREATE INDEX idx_user_locations_device_id ON user_locations (device_id) WHERE device_id IS NOT NULL;

-- User Location History: Migrated to Redis (see docs/REDIS_LOCATION_HISTORY.md)
-- Table removed in migration 000003_drop_location_history_table.up.sql