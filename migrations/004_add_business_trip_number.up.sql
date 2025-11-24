-- Add business_trip_number column
ALTER TABLE business_trips ADD COLUMN business_trip_number VARCHAR(10) NULL;

-- Create a unique index for business trip number
CREATE UNIQUE INDEX idx_business_trips_business_trip_number ON business_trips(business_trip_number) WHERE business_trip_number IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN business_trips.business_trip_number IS 'Auto-generated business trip number in format BT-XXXX';