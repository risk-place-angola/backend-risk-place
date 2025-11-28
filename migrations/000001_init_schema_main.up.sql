-- ============================================================================
-- RISKPLACE DATABASE SCHEMA - PRODUCTION DEPLOYMENT
-- Generated from development database: 2025-11-27
-- ============================================================================
-- This schema represents the complete and verified structure from the 
-- development environment, ready for production deployment.
-- ============================================================================

-- ============================================================================
-- EXTENSIONS
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS cube WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS earthdistance WITH SCHEMA public;

-- ============================================================================
-- CUSTOM TYPES (ENUMS)
-- ============================================================================

CREATE TYPE public.alert_severity AS ENUM (
    'low',
    'medium',
    'high',
    'critical'
);

CREATE TYPE public.alert_status AS ENUM (
    'active',
    'resolved',
    'expired'
);

CREATE TYPE public.entity_type AS ENUM (
    'erce',
    'erfce'
);

CREATE TYPE public.notification_type AS ENUM (
    'alert',
    'report'
);

CREATE TYPE public.report_status AS ENUM (
    'pending',
    'verified',
    'resolved',
    'rejected'
);

CREATE TYPE public.role_type AS ENUM (
    'citizen',
    'erce',
    'erfce',
    'admin'
);

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Function: Auto-expire anonymous alerts after 2 hours
CREATE FUNCTION public.auto_expire_anonymous_alerts() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE alerts
    SET status = 'expired'
    WHERE anonymous_session_id IS NOT NULL
      AND status = 'active'
      AND created_at < NOW() - INTERVAL '2 hours';
END;
$$;

-- Function: Calculate distance in meters between two coordinates
CREATE FUNCTION public.calculate_distance_meters(
    lat1 double precision, 
    lon1 double precision, 
    lat2 double precision, 
    lon2 double precision
) RETURNS integer
    LANGUAGE plpgsql IMMUTABLE
    AS $$
BEGIN
    RETURN CAST(earth_distance(
        ll_to_earth(lat1, lon1),
        ll_to_earth(lat2, lon2)
    ) AS INTEGER);
END;
$$;

-- Function: Cleanup old anonymous sessions (30+ days)
CREATE FUNCTION public.cleanup_old_anonymous_sessions() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM anonymous_sessions
    WHERE last_seen < NOW() - INTERVAL '30 days';
END;
$$;

-- Function: Trigger to set updated_at timestamp
CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- Function: Update anonymous session last_seen on activity
CREATE FUNCTION public.update_anonymous_session_last_seen() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE anonymous_sessions
    SET last_seen = NOW()
    WHERE id = NEW.anonymous_session_id;
    RETURN NEW;
END;
$$;

-- ============================================================================
-- TABLES
-- ============================================================================

-- Table: users
CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    phone text,
    latitude double precision,
    longitude double precision,
    alert_radius_meters integer DEFAULT 1000,
    email_verified boolean DEFAULT false,
    email_verification_code text,
    email_verification_expires_at timestamp without time zone,
    nif text,
    province text,
    municipality text,
    neighborhood text,
    address text,
    zip_code text,
    country text,
    last_login timestamp without time zone,
    failed_attempts integer DEFAULT 0,
    locked_until timestamp without time zone,
    device_fcm_token text,
    device_language text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone,
    home_address_name character varying(255),
    home_address_address text,
    home_address_lat double precision,
    home_address_lon double precision,
    work_address_name character varying(255),
    work_address_address text,
    work_address_lat double precision,
    work_address_lon double precision,
    linked_device_id text,
    push_notification_enabled boolean DEFAULT true,
    sms_notification_enabled boolean DEFAULT false,
    trust_score integer DEFAULT 50,
    reports_submitted integer DEFAULT 0,
    reports_verified integer DEFAULT 0,
    account_verified boolean DEFAULT false
);

-- Table: roles
CREATE TABLE public.roles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    description text
);

-- Table: permissions
CREATE TABLE public.permissions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    resource text NOT NULL,
    action text NOT NULL,
    code text GENERATED ALWAYS AS (lower(((resource || ':'::text) || action))) STORED
);

-- Table: role_permissions
CREATE TABLE public.role_permissions (
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    granted_at timestamp without time zone DEFAULT now() NOT NULL
);

-- Table: user_roles
CREATE TABLE public.user_roles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    assigned_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

-- Table: anonymous_sessions
CREATE TABLE public.anonymous_sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    device_id text NOT NULL,
    device_fcm_token text,
    device_platform text,
    device_model text,
    latitude double precision,
    longitude double precision,
    alert_radius_meters integer DEFAULT 1000,
    device_language text DEFAULT 'pt'::text,
    last_seen timestamp without time zone DEFAULT now(),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    migrated_to_user_id uuid,
    migrated_at timestamp without time zone,
    is_active boolean DEFAULT true,
    push_notification_enabled boolean DEFAULT true,
    sms_notification_enabled boolean DEFAULT false
);

-- Table: emergency_contacts
CREATE TABLE public.emergency_contacts (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    phone character varying(20) NOT NULL,
    relation character varying(50) NOT NULL,
    is_priority boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT emergency_contacts_relation_check CHECK (((relation)::text = ANY ((ARRAY['family'::character varying, 'friend'::character varying, 'colleague'::character varying, 'neighbor'::character varying, 'other'::character varying])::text[])))
);

-- Table: location_sharings
CREATE TABLE public.location_sharings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid,
    anonymous_session_id uuid,
    device_id text,
    owner_name character varying(255),
    token character varying(255) NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    duration_minutes integer NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    last_updated_at timestamp without time zone DEFAULT now() NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT location_sharing_owner_check CHECK ((((user_id IS NOT NULL) AND (anonymous_session_id IS NULL) AND (device_id IS NULL)) OR ((user_id IS NULL) AND (anonymous_session_id IS NOT NULL) AND (device_id IS NOT NULL))))
);

-- Table: risk_types
CREATE TABLE public.risk_types (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    default_radius_meters integer DEFAULT 500,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    icon_path text,
    is_sensitive boolean DEFAULT false,
    is_enabled boolean DEFAULT true NOT NULL
);

COMMENT ON COLUMN public.risk_types.is_enabled IS 'Controls visibility of risk types in mobile app. When false, all associated reports are also hidden.';

-- Table: risk_topics
CREATE TABLE public.risk_topics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    risk_type_id uuid NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    icon_path text,
    is_sensitive boolean DEFAULT false NOT NULL
);

-- Table: reports
CREATE TABLE public.reports (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    risk_type_id uuid NOT NULL,
    risk_topic_id uuid,
    description text,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    province text,
    municipality text,
    neighborhood text,
    address text,
    image_url text,
    status public.report_status DEFAULT 'pending'::public.report_status,
    reviewed_by uuid,
    resolved_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    verification_count integer DEFAULT 0,
    rejection_count integer DEFAULT 0,
    expires_at timestamp without time zone,
    is_private boolean DEFAULT false NOT NULL
);

-- Table: alerts
CREATE TABLE public.alerts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_by uuid,
    risk_type_id uuid NOT NULL,
    risk_topic_id uuid,
    message text NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    province character varying(100),
    municipality character varying(100),
    neighborhood character varying(100),
    address text,
    radius_meters integer DEFAULT 500 NOT NULL,
    severity public.alert_severity DEFAULT 'medium'::public.alert_severity,
    status public.alert_status DEFAULT 'active'::public.alert_status,
    created_at timestamp without time zone DEFAULT now(),
    expires_at timestamp without time zone,
    resolved_at timestamp without time zone,
    anonymous_session_id uuid,
    device_id text,
    CONSTRAINT alerts_creator_check CHECK ((((created_by IS NOT NULL) AND (anonymous_session_id IS NULL) AND (device_id IS NULL)) OR ((created_by IS NULL) AND (anonymous_session_id IS NOT NULL) AND (device_id IS NOT NULL))))
);

-- Table: alert_subscriptions
CREATE TABLE public.alert_subscriptions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    alert_id uuid NOT NULL,
    user_id uuid,
    subscribed_at timestamp without time zone DEFAULT now() NOT NULL,
    anonymous_session_id uuid,
    device_id text,
    CONSTRAINT alert_subscriptions_subscriber_check CHECK ((((user_id IS NOT NULL) AND (anonymous_session_id IS NULL) AND (device_id IS NULL)) OR ((user_id IS NULL) AND (anonymous_session_id IS NOT NULL) AND (device_id IS NOT NULL))))
);

-- Table: report_votes
CREATE TABLE public.report_votes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    report_id uuid NOT NULL,
    user_id uuid,
    anonymous_session_id uuid,
    vote_type character varying(10) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT report_votes_vote_type_check CHECK (((vote_type)::text = ANY ((ARRAY['upvote'::character varying, 'downvote'::character varying])::text[]))),
    CONSTRAINT vote_user_or_anonymous CHECK ((((user_id IS NOT NULL) AND (anonymous_session_id IS NULL)) OR ((user_id IS NULL) AND (anonymous_session_id IS NOT NULL))))
);

-- Table: notifications
CREATE TABLE public.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    type public.notification_type NOT NULL,
    reference_id uuid NOT NULL,
    user_id uuid NOT NULL,
    sent_at timestamp without time zone DEFAULT now(),
    seen_at timestamp without time zone
);

-- Table: entities
CREATE TABLE public.entities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    entity_type public.entity_type NOT NULL,
    province text,
    municipality text,
    latitude double precision,
    longitude double precision,
    contact_email text,
    contact_phone text,
    created_at timestamp without time zone DEFAULT now()
);

-- Table: user_safety_settings
CREATE TABLE public.user_safety_settings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    notifications_enabled boolean DEFAULT true,
    notification_alert_types text[] DEFAULT ARRAY['high'::text, 'critical'::text],
    notification_alert_radius_mins integer DEFAULT 1000,
    notification_report_types text[] DEFAULT ARRAY['verified'::text],
    notification_report_radius_mins integer DEFAULT 500,
    location_sharing_enabled boolean DEFAULT false,
    location_history_enabled boolean DEFAULT true,
    profile_visibility text DEFAULT 'public'::text,
    anonymous_reports boolean DEFAULT false,
    show_online_status boolean DEFAULT true,
    auto_alerts_enabled boolean DEFAULT false,
    danger_zones_enabled boolean DEFAULT true,
    time_based_alerts_enabled boolean DEFAULT false,
    high_risk_start_time time without time zone DEFAULT '22:00:00'::time without time zone,
    high_risk_end_time time without time zone DEFAULT '06:00:00'::time without time zone,
    night_mode_enabled boolean DEFAULT false,
    night_mode_start_time time without time zone DEFAULT '22:00:00'::time without time zone,
    night_mode_end_time time without time zone DEFAULT '06:00:00'::time without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    anonymous_session_id uuid,
    device_id text,
    CONSTRAINT user_safety_settings_owner_check CHECK ((((user_id IS NOT NULL) AND (anonymous_session_id IS NULL) AND (device_id IS NULL)) OR ((user_id IS NULL) AND (anonymous_session_id IS NOT NULL) AND (device_id IS NOT NULL)))),
    CONSTRAINT user_safety_settings_profile_visibility_check CHECK ((profile_visibility = ANY (ARRAY['public'::text, 'friends'::text, 'private'::text])))
);

-- Table: device_user_mappings
CREATE TABLE public.device_user_mappings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    device_id text NOT NULL,
    anonymous_session_id uuid NOT NULL,
    user_id uuid NOT NULL,
    mapped_at timestamp without time zone DEFAULT now() NOT NULL,
    unmapped_at timestamp without time zone,
    is_active boolean DEFAULT true NOT NULL,
    CONSTRAINT device_user_mapping_active_check CHECK ((((is_active = true) AND (unmapped_at IS NULL)) OR ((is_active = false) AND (unmapped_at IS NOT NULL))))
);

-- Table: anonymous_user_migrations
CREATE TABLE public.anonymous_user_migrations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    anonymous_session_id uuid NOT NULL,
    device_id text NOT NULL,
    user_id uuid NOT NULL,
    alerts_migrated integer DEFAULT 0 NOT NULL,
    subscriptions_migrated integer DEFAULT 0 NOT NULL,
    settings_migrated boolean DEFAULT false NOT NULL,
    location_sharings_migrated integer DEFAULT 0 NOT NULL,
    migration_type character varying(20) NOT NULL,
    started_at timestamp without time zone DEFAULT now() NOT NULL,
    completed_at timestamp without time zone,
    failed_at timestamp without time zone,
    error_message text,
    CONSTRAINT anonymous_user_migrations_migration_type_check CHECK (((migration_type)::text = ANY ((ARRAY['signup'::character varying, 'login'::character varying, 'manual'::character varying])::text[]))),
    CONSTRAINT migration_completion_check CHECK ((((completed_at IS NOT NULL) AND (failed_at IS NULL)) OR ((completed_at IS NULL) AND (failed_at IS NOT NULL)) OR ((completed_at IS NULL) AND (failed_at IS NULL))))
);

-- Table: user_locations (Waze-style nearby users feature)
-- Uses earthdistance extension for efficient geospatial queries
CREATE TABLE public.user_locations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    device_id character varying(255),
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    speed double precision DEFAULT 0,
    heading double precision DEFAULT 0,
    avatar_id integer NOT NULL,
    color character varying(7) NOT NULL,
    is_anonymous boolean DEFAULT false,
    last_update timestamp with time zone DEFAULT now(),
    created_at timestamp with time zone DEFAULT now()
);

-- ============================================================================
-- PRIMARY KEYS
-- ============================================================================

ALTER TABLE ONLY public.users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.roles ADD CONSTRAINT roles_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.permissions ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.role_permissions ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id);
ALTER TABLE ONLY public.user_roles ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.anonymous_sessions ADD CONSTRAINT anonymous_sessions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.emergency_contacts ADD CONSTRAINT emergency_contacts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.location_sharings ADD CONSTRAINT location_sharings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.risk_types ADD CONSTRAINT risk_types_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.risk_topics ADD CONSTRAINT risk_topics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.reports ADD CONSTRAINT reports_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.alerts ADD CONSTRAINT alerts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.alert_subscriptions ADD CONSTRAINT alert_subscriptions_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.report_votes ADD CONSTRAINT report_votes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.notifications ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.entities ADD CONSTRAINT entities_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_safety_settings ADD CONSTRAINT user_safety_settings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.device_user_mappings ADD CONSTRAINT device_user_mappings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.anonymous_user_migrations ADD CONSTRAINT anonymous_user_migrations_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_locations ADD CONSTRAINT user_locations_pkey PRIMARY KEY (id);

-- ============================================================================
-- UNIQUE CONSTRAINTS
-- ============================================================================

ALTER TABLE ONLY public.users ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users ADD CONSTRAINT users_nif_key UNIQUE (nif);
ALTER TABLE ONLY public.roles ADD CONSTRAINT roles_name_key UNIQUE (name);
ALTER TABLE ONLY public.permissions ADD CONSTRAINT permissions_resource_action_key UNIQUE (resource, action);
ALTER TABLE ONLY public.user_roles ADD CONSTRAINT user_roles_user_id_role_id_key UNIQUE (user_id, role_id);
ALTER TABLE ONLY public.anonymous_sessions ADD CONSTRAINT anonymous_sessions_device_id_key UNIQUE (device_id);
ALTER TABLE ONLY public.location_sharings ADD CONSTRAINT location_sharings_token_key UNIQUE (token);
ALTER TABLE ONLY public.risk_types ADD CONSTRAINT risk_types_name_key UNIQUE (name);
ALTER TABLE ONLY public.risk_topics ADD CONSTRAINT risk_topics_risk_type_id_name_key UNIQUE (risk_type_id, name);
ALTER TABLE ONLY public.report_votes ADD CONSTRAINT report_votes_report_id_user_id_key UNIQUE (report_id, user_id);
ALTER TABLE ONLY public.report_votes ADD CONSTRAINT report_votes_report_id_anonymous_session_id_key UNIQUE (report_id, anonymous_session_id);
ALTER TABLE ONLY public.notifications ADD CONSTRAINT notifications_type_reference_id_user_id_key UNIQUE (type, reference_id, user_id);
ALTER TABLE ONLY public.user_locations ADD CONSTRAINT user_locations_user_id_key UNIQUE (user_id);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Users indexes
CREATE INDEX idx_users_email_trgm ON public.users USING gin (email public.gin_trgm_ops);
CREATE INDEX idx_users_geo ON public.users USING btree (latitude, longitude);
CREATE INDEX idx_users_trust_score ON public.users USING btree (trust_score);
CREATE INDEX idx_users_notification_prefs ON public.users USING btree (push_notification_enabled, sms_notification_enabled) WHERE ((push_notification_enabled = true) OR (sms_notification_enabled = true));
CREATE INDEX idx_users_home_location ON public.users USING gist (public.ll_to_earth(home_address_lat, home_address_lon)) WHERE ((home_address_lat IS NOT NULL) AND (home_address_lon IS NOT NULL));
CREATE INDEX idx_users_work_location ON public.users USING gist (public.ll_to_earth(work_address_lat, work_address_lon)) WHERE ((work_address_lat IS NOT NULL) AND (work_address_lon IS NOT NULL));
CREATE INDEX idx_users_linked_device_id ON public.users USING btree (linked_device_id) WHERE (linked_device_id IS NOT NULL);

-- Roles and permissions indexes
CREATE INDEX idx_roles_priority ON public.roles USING btree (priority DESC);
CREATE INDEX idx_user_roles_user ON public.user_roles USING btree (user_id);
CREATE INDEX idx_role_permissions_role ON public.role_permissions USING btree (role_id);

-- Anonymous sessions indexes
CREATE INDEX idx_anonymous_sessions_device_id ON public.anonymous_sessions USING btree (device_id);
CREATE INDEX idx_anonymous_sessions_location ON public.anonymous_sessions USING btree (latitude, longitude);
CREATE INDEX idx_anonymous_sessions_last_seen ON public.anonymous_sessions USING btree (last_seen);
CREATE INDEX idx_anonymous_sessions_is_active ON public.anonymous_sessions USING btree (is_active) WHERE (is_active = true);
CREATE INDEX idx_anonymous_sessions_migrated_user ON public.anonymous_sessions USING btree (migrated_to_user_id) WHERE (migrated_to_user_id IS NOT NULL);
CREATE INDEX idx_anonymous_sessions_notification_prefs ON public.anonymous_sessions USING btree (push_notification_enabled, sms_notification_enabled) WHERE ((push_notification_enabled = true) OR (sms_notification_enabled = true));

-- Emergency contacts indexes
CREATE INDEX idx_emergency_contacts_user_id ON public.emergency_contacts USING btree (user_id);
CREATE INDEX idx_emergency_contacts_priority ON public.emergency_contacts USING btree (user_id, is_priority) WHERE (is_priority = true);

-- Location sharings indexes
CREATE INDEX idx_location_sharings_token ON public.location_sharings USING btree (token);
CREATE INDEX idx_location_sharings_user_id ON public.location_sharings USING btree (user_id);
CREATE INDEX idx_location_sharings_anonymous_session_id ON public.location_sharings USING btree (anonymous_session_id);
CREATE INDEX idx_location_sharings_device_id ON public.location_sharings USING btree (device_id);
CREATE INDEX idx_location_sharings_expires_at ON public.location_sharings USING btree (expires_at);
CREATE INDEX idx_location_sharings_is_active ON public.location_sharings USING btree (is_active);

-- Risk types indexes
CREATE INDEX idx_risk_types_is_enabled ON public.risk_types USING btree (is_enabled) WHERE (is_enabled = true);

-- Reports indexes
CREATE INDEX idx_reports_user ON public.reports USING btree (user_id);
CREATE INDEX idx_reports_status ON public.reports USING btree (status);
CREATE INDEX idx_reports_created_at ON public.reports USING btree (created_at DESC);
CREATE INDEX idx_reports_status_created ON public.reports USING btree (status, created_at);
CREATE INDEX idx_reports_expires_at ON public.reports USING btree (expires_at);
CREATE INDEX idx_reports_is_private ON public.reports USING btree (is_private);
CREATE INDEX idx_reports_status_private ON public.reports USING btree (status, is_private);
CREATE INDEX idx_reports_location_gist ON public.reports USING gist (public.ll_to_earth(latitude, longitude));

-- Alerts indexes
CREATE INDEX idx_alerts_location_gist ON public.alerts USING gist (public.ll_to_earth(latitude, longitude));
CREATE INDEX idx_alerts_anonymous_session_id ON public.alerts USING btree (anonymous_session_id) WHERE (anonymous_session_id IS NOT NULL);
CREATE INDEX idx_alerts_device_id ON public.alerts USING btree (device_id) WHERE (device_id IS NOT NULL);

-- Alert subscriptions indexes
CREATE INDEX idx_alert_subscriptions_alert_id ON public.alert_subscriptions USING btree (alert_id);
CREATE INDEX idx_alert_subscriptions_user_id ON public.alert_subscriptions USING btree (user_id);
CREATE INDEX idx_alert_subscriptions_subscribed_at ON public.alert_subscriptions USING btree (subscribed_at DESC);
CREATE INDEX idx_alert_subscriptions_anonymous_session_id ON public.alert_subscriptions USING btree (anonymous_session_id) WHERE (anonymous_session_id IS NOT NULL);
CREATE INDEX idx_alert_subscriptions_device_id ON public.alert_subscriptions USING btree (device_id) WHERE (device_id IS NOT NULL);
CREATE UNIQUE INDEX idx_alert_subscriptions_unique_user ON public.alert_subscriptions USING btree (alert_id, user_id) WHERE (user_id IS NOT NULL);
CREATE UNIQUE INDEX idx_alert_subscriptions_unique_device ON public.alert_subscriptions USING btree (alert_id, device_id) WHERE (device_id IS NOT NULL);

-- Report votes indexes
CREATE INDEX idx_report_votes_report ON public.report_votes USING btree (report_id);
CREATE INDEX idx_report_votes_user ON public.report_votes USING btree (user_id);

-- User safety settings indexes
CREATE INDEX idx_user_safety_settings_user_id ON public.user_safety_settings USING btree (user_id);
CREATE INDEX idx_user_safety_settings_updated_at ON public.user_safety_settings USING btree (updated_at DESC);
CREATE INDEX idx_safety_settings_anonymous_session_id ON public.user_safety_settings USING btree (anonymous_session_id) WHERE (anonymous_session_id IS NOT NULL);
CREATE INDEX idx_safety_settings_device_id ON public.user_safety_settings USING btree (device_id) WHERE (device_id IS NOT NULL);
CREATE INDEX idx_safety_settings_notifications_enabled ON public.user_safety_settings USING btree (notifications_enabled) WHERE (notifications_enabled = true);
CREATE UNIQUE INDEX idx_safety_settings_unique_user ON public.user_safety_settings USING btree (user_id) WHERE (user_id IS NOT NULL);
CREATE UNIQUE INDEX idx_safety_settings_unique_device ON public.user_safety_settings USING btree (device_id) WHERE (device_id IS NOT NULL);

-- Device user mappings indexes
CREATE UNIQUE INDEX idx_device_user_mappings_active ON public.device_user_mappings USING btree (device_id, user_id) WHERE (is_active = true);
CREATE INDEX idx_device_user_mappings_device_id ON public.device_user_mappings USING btree (device_id);
CREATE INDEX idx_device_user_mappings_user_id ON public.device_user_mappings USING btree (user_id);
CREATE INDEX idx_device_user_mappings_anonymous_session_id ON public.device_user_mappings USING btree (anonymous_session_id);

-- Anonymous user migrations indexes
CREATE INDEX idx_migrations_anonymous_session_id ON public.anonymous_user_migrations USING btree (anonymous_session_id);
CREATE INDEX idx_migrations_user_id ON public.anonymous_user_migrations USING btree (user_id);
CREATE INDEX idx_migrations_device_id ON public.anonymous_user_migrations USING btree (device_id);
CREATE INDEX idx_migrations_started_at ON public.anonymous_user_migrations USING btree (started_at DESC);

-- User locations indexes
CREATE INDEX idx_user_locations_user_id ON public.user_locations USING btree (user_id);
CREATE INDEX idx_user_locations_last_update ON public.user_locations USING btree (last_update);
CREATE INDEX idx_user_locations_location ON public.user_locations USING gist (public.ll_to_earth(latitude, longitude));
CREATE INDEX idx_user_locations_device_id ON public.user_locations USING btree (device_id) WHERE (device_id IS NOT NULL);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Trigger: Update timestamp on users table
CREATE TRIGGER set_timestamp 
    BEFORE UPDATE ON public.users 
    FOR EACH ROW 
    EXECUTE FUNCTION public.trigger_set_timestamp();

-- Trigger: Update timestamp on reports table
CREATE TRIGGER set_timestamp_reports 
    BEFORE UPDATE ON public.reports 
    FOR EACH ROW 
    EXECUTE FUNCTION public.trigger_set_timestamp();

-- Trigger: Update timestamp on anonymous_sessions table
CREATE TRIGGER set_timestamp_anonymous_sessions 
    BEFORE UPDATE ON public.anonymous_sessions 
    FOR EACH ROW 
    EXECUTE FUNCTION public.trigger_set_timestamp();

-- Trigger: Update last_seen on alert creation/update
CREATE TRIGGER trigger_update_last_seen_on_alert 
    AFTER INSERT OR UPDATE ON public.alerts 
    FOR EACH ROW 
    WHEN ((new.anonymous_session_id IS NOT NULL)) 
    EXECUTE FUNCTION public.update_anonymous_session_last_seen();

-- Trigger: Update last_seen on subscription creation/update
CREATE TRIGGER trigger_update_last_seen_on_subscription 
    AFTER INSERT OR UPDATE ON public.alert_subscriptions 
    FOR EACH ROW 
    WHEN ((new.anonymous_session_id IS NOT NULL)) 
    EXECUTE FUNCTION public.update_anonymous_session_last_seen();

-- ============================================================================
-- FOREIGN KEY CONSTRAINTS
-- ============================================================================

-- User roles foreign keys
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;

-- Role permissions foreign keys
ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;

-- Anonymous sessions foreign keys
ALTER TABLE ONLY public.anonymous_sessions
    ADD CONSTRAINT anonymous_sessions_migrated_to_user_id_fkey FOREIGN KEY (migrated_to_user_id) REFERENCES public.users(id) ON DELETE SET NULL;

-- Emergency contacts foreign keys
ALTER TABLE ONLY public.emergency_contacts
    ADD CONSTRAINT emergency_contacts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

-- Location sharings foreign keys
ALTER TABLE ONLY public.location_sharings
    ADD CONSTRAINT location_sharings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.location_sharings
    ADD CONSTRAINT location_sharings_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE CASCADE;

-- Risk topics foreign keys
ALTER TABLE ONLY public.risk_topics
    ADD CONSTRAINT risk_topics_risk_type_id_fkey FOREIGN KEY (risk_type_id) REFERENCES public.risk_types(id) ON DELETE CASCADE;

-- Reports foreign keys
ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_risk_type_id_fkey FOREIGN KEY (risk_type_id) REFERENCES public.risk_types(id);
ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_risk_topic_id_fkey FOREIGN KEY (risk_topic_id) REFERENCES public.risk_topics(id);
ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES public.users(id);

-- Alerts foreign keys
ALTER TABLE ONLY public.alerts
    ADD CONSTRAINT alerts_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.alerts
    ADD CONSTRAINT alerts_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.alerts
    ADD CONSTRAINT alerts_risk_type_id_fkey FOREIGN KEY (risk_type_id) REFERENCES public.risk_types(id);
ALTER TABLE ONLY public.alerts
    ADD CONSTRAINT alerts_risk_topic_id_fkey FOREIGN KEY (risk_topic_id) REFERENCES public.risk_topics(id);

-- Alert subscriptions foreign keys
ALTER TABLE ONLY public.alert_subscriptions
    ADD CONSTRAINT alert_subscriptions_alert_id_fkey FOREIGN KEY (alert_id) REFERENCES public.alerts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.alert_subscriptions
    ADD CONSTRAINT alert_subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.alert_subscriptions
    ADD CONSTRAINT alert_subscriptions_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE CASCADE;

-- Report votes foreign keys
ALTER TABLE ONLY public.report_votes
    ADD CONSTRAINT report_votes_report_id_fkey FOREIGN KEY (report_id) REFERENCES public.reports(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.report_votes
    ADD CONSTRAINT report_votes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.report_votes
    ADD CONSTRAINT report_votes_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE CASCADE;

-- Notifications foreign keys
ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

-- User safety settings foreign keys
ALTER TABLE ONLY public.user_safety_settings
    ADD CONSTRAINT user_safety_settings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.user_safety_settings
    ADD CONSTRAINT user_safety_settings_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE CASCADE;

-- Device user mappings foreign keys
ALTER TABLE ONLY public.device_user_mappings
    ADD CONSTRAINT device_user_mappings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.device_user_mappings
    ADD CONSTRAINT device_user_mappings_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE CASCADE;

-- Anonymous user migrations foreign keys
ALTER TABLE ONLY public.anonymous_user_migrations
    ADD CONSTRAINT anonymous_user_migrations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.anonymous_user_migrations
    ADD CONSTRAINT anonymous_user_migrations_anonymous_session_id_fkey FOREIGN KEY (anonymous_session_id) REFERENCES public.anonymous_sessions(id) ON DELETE CASCADE;

-- ============================================================================
-- DEPLOYMENT COMPLETE
-- ============================================================================
-- Schema successfully created from development database
-- All tables, indexes, constraints, triggers, and functions are in place
-- Ready for production use
-- ============================================================================
