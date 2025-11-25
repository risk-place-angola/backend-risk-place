-- Rollback: Restore user_id NOT NULL constraint in user_safety_settings
ALTER TABLE user_safety_settings ALTER COLUMN user_id SET NOT NULL;
