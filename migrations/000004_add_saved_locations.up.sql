ALTER TABLE users
ADD COLUMN home_address_name VARCHAR(255),
ADD COLUMN home_address_address TEXT,
ADD COLUMN home_address_lat DOUBLE PRECISION,
ADD COLUMN home_address_lon DOUBLE PRECISION,
ADD COLUMN work_address_name VARCHAR(255),
ADD COLUMN work_address_address TEXT,
ADD COLUMN work_address_lat DOUBLE PRECISION,
ADD COLUMN work_address_lon DOUBLE PRECISION;

CREATE INDEX IF NOT EXISTS idx_users_home_location ON users USING GIST (
    ST_SetSRID(ST_MakePoint(home_address_lon, home_address_lat), 4326)
) WHERE home_address_lat IS NOT NULL AND home_address_lon IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_work_location ON users USING GIST (
    ST_SetSRID(ST_MakePoint(work_address_lon, work_address_lat), 4326)
) WHERE work_address_lat IS NOT NULL AND work_address_lon IS NOT NULL;
