-- Make user_id nullable in alert_subscriptions to support anonymous subscriptions
ALTER TABLE alert_subscriptions ALTER COLUMN user_id DROP NOT NULL;
