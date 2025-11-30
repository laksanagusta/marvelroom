-- Migration: Update business trip status constraint to include ready_to_verify
-- Description: Updates the status check constraint to include the new ready_to_verify status

-- Drop existing constraint
ALTER TABLE business_trips DROP CONSTRAINT IF EXISTS chk_business_trip_status;

-- Add updated constraint with ready_to_verify status
ALTER TABLE business_trips ADD CONSTRAINT chk_business_trip_status
    CHECK (status IN ('draft', 'ready_to_verify', 'ongoing', 'completed', 'canceled'));

-- Update comment to include new status
COMMENT ON COLUMN business_trips.status IS 'Status of the business trip: draft, ready_to_verify, ongoing, completed, or canceled';