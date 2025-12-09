-- Migration: Create business trip histories table
-- Description: Adds support for tracking all important changes on business trips (status changes and verification process)

-- Create business_trip_histories table
CREATE TABLE IF NOT EXISTS business_trip_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_trip_id UUID NOT NULL REFERENCES business_trips(id) ON DELETE CASCADE,
    change_type VARCHAR(50) NOT NULL,
    field_name VARCHAR(100) NULL,
    old_value TEXT NULL,
    new_value TEXT NULL,
    changed_by_user_id VARCHAR(100) NULL,
    changed_by_user_name VARCHAR(255) NULL,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_history_change_type CHECK (change_type IN ('status_change', 'verification_approved', 'verification_rejected', 'verification_pending'))
);

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_business_trip_histories_business_trip_id ON business_trip_histories(business_trip_id);
CREATE INDEX IF NOT EXISTS idx_business_trip_histories_change_type ON business_trip_histories(change_type);
CREATE INDEX IF NOT EXISTS idx_business_trip_histories_created_at ON business_trip_histories(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_business_trip_histories_changed_by_user_id ON business_trip_histories(changed_by_user_id);

-- Add comments for documentation
COMMENT ON TABLE business_trip_histories IS 'Table for tracking all important changes on business trips (status changes and verification process)';
COMMENT ON COLUMN business_trip_histories.business_trip_id IS 'Reference to the business trip';
COMMENT ON COLUMN business_trip_histories.change_type IS 'Type of change (status_change, verification_approved, verification_rejected, verification_pending)';
COMMENT ON COLUMN business_trip_histories.field_name IS 'Name of the field that was changed (e.g., "status", "verificator")';
COMMENT ON COLUMN business_trip_histories.old_value IS 'Previous value before the change';
COMMENT ON COLUMN business_trip_histories.new_value IS 'New value after the change';
COMMENT ON COLUMN business_trip_histories.changed_by_user_id IS 'User ID who made the change';
COMMENT ON COLUMN business_trip_histories.changed_by_user_name IS 'Full name of the user who made the change';
COMMENT ON COLUMN business_trip_histories.notes IS 'Additional notes or comments about the change';
COMMENT ON COLUMN business_trip_histories.created_at IS 'Timestamp when the change was recorded';
