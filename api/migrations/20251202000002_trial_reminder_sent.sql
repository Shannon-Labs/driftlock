-- +goose Up
-- Add trial_reminder_sent flag for idempotent trial ending notifications
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS trial_reminder_sent BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE tenants DROP COLUMN IF EXISTS trial_reminder_sent;
