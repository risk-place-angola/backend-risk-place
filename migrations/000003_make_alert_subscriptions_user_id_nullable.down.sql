-- Rollback: Restore user_id NOT NULL constraint in alert_subscriptions
ALTER TABLE alert_subscriptions ALTER COLUMN user_id SET NOT NULL;
