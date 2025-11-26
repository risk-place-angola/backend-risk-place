DROP FUNCTION IF EXISTS calculate_distance_meters(DOUBLE PRECISION, DOUBLE PRECISION, DOUBLE PRECISION, DOUBLE PRECISION);
DROP INDEX IF EXISTS idx_reports_location_gist;
DROP INDEX IF EXISTS idx_alerts_location_gist;
DROP INDEX IF EXISTS idx_safety_settings_notifications_enabled;

-- Note: Extensions are NOT dropped as they may be used by other parts of the database
-- If you need to drop them manually, run:
-- DROP EXTENSION IF EXISTS earthdistance;
-- DROP EXTENSION IF EXISTS cube;
