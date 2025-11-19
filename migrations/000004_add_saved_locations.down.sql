DROP INDEX IF EXISTS idx_users_home_location;
DROP INDEX IF EXISTS idx_users_work_location;

ALTER TABLE users
DROP COLUMN IF EXISTS home_address_name,
DROP COLUMN IF EXISTS home_address_address,
DROP COLUMN IF EXISTS home_address_lat,
DROP COLUMN IF EXISTS home_address_lon,
DROP COLUMN IF EXISTS work_address_name,
DROP COLUMN IF EXISTS work_address_address,
DROP COLUMN IF EXISTS work_address_lat,
DROP COLUMN IF EXISTS work_address_lon;
