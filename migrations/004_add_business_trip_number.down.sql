-- Remove the unique index
DROP INDEX IF EXISTS idx_business_trips_business_trip_number;

-- Remove the business_trip_number column
ALTER TABLE business_trips DROP COLUMN IF EXISTS business_trip_number;