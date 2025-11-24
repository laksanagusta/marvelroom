-- Rollback: Remove status and document_link columns from business_trips table
-- Description: Removes status tracking and document link functionality

-- Drop the constraint first
ALTER TABLE business_trips DROP CONSTRAINT IF EXISTS chk_business_trip_status;

-- Drop the index
DROP INDEX IF EXISTS idx_business_trips_status;

-- Drop the columns
ALTER TABLE business_trips DROP COLUMN IF EXISTS status;
ALTER TABLE business_trips DROP COLUMN IF EXISTS document_link;