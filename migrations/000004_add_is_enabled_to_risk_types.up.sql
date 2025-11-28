-- Add is_enabled column to risk_types table
-- This allows admins to enable/disable specific risk types
-- When disabled, risk types and their associated reports won't be visible in the mobile app

ALTER TABLE risk_types 
ADD COLUMN IF NOT EXISTS is_enabled BOOLEAN NOT NULL DEFAULT TRUE;

-- Create an index for better query performance when filtering by enabled risk types
CREATE INDEX idx_risk_types_is_enabled ON risk_types(is_enabled) WHERE is_enabled = TRUE;

-- Add comment to explain the column purpose
COMMENT ON COLUMN risk_types.is_enabled IS 'Controls visibility of risk types in mobile app. When false, all associated reports are also hidden.';
