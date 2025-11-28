-- Make user_id nullable in user_safety_settings to support anonymous users
ALTER TABLE user_safety_settings ALTER COLUMN user_id DROP NOT NULL;
