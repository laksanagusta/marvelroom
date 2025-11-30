-- Migration: Revert business trip status constraint to original values
-- Description: Removes ready_to_verify status from the check constraint

-- Drop existing constraint
ALTER TABLE business_trips DROP CONSTRAINT IF EXISTS chk_business_trip_status;

-- Add original constraint without ready_to_verify status
ALTER TABLE business_trips ADD CONSTRAINT chk_business_trip_status
    CHECK (status IN ('draft', 'ongoing', 'completed', 'canceled'));

-- Revert comment to original
COMMENT ON COLUMN business_trips.status IS 'Status of the business trip: draft, ongoing, completed, or canceled';