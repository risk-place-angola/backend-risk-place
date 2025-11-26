-- Enable required PostgreSQL extensions for geospatial queries
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

-- Add index for notification filtering
CREATE INDEX IF NOT EXISTS idx_safety_settings_notifications_enabled 
ON user_safety_settings(notifications_enabled) 
WHERE notifications_enabled = true;

-- Add GiST index for geospatial queries on alerts
CREATE INDEX IF NOT EXISTS idx_alerts_location_gist 
ON alerts USING GIST (ll_to_earth(latitude, longitude));

-- Add GiST index for geospatial queries on reports  
CREATE INDEX IF NOT EXISTS idx_reports_location_gist 
ON reports USING GIST (ll_to_earth(latitude, longitude));

-- Function to calculate distance in meters
CREATE OR REPLACE FUNCTION calculate_distance_meters(
    lat1 DOUBLE PRECISION,
    lon1 DOUBLE PRECISION, 
    lat2 DOUBLE PRECISION,
    lon2 DOUBLE PRECISION
) RETURNS INTEGER AS $$
BEGIN
    RETURN CAST(earth_distance(
        ll_to_earth(lat1, lon1),
        ll_to_earth(lat2, lon2)
    ) AS INTEGER);
END;
$$ LANGUAGE plpgsql IMMUTABLE;
