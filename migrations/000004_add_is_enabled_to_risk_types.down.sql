-- Rollback migration: Remove is_enabled column from risk_types table

DROP INDEX IF EXISTS idx_risk_types_is_enabled;

ALTER TABLE risk_types 
DROP COLUMN IF EXISTS is_enabled;
