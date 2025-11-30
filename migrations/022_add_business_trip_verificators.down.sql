-- Migration: Remove verificators from business trips
-- Description: Removes the business_trip_verificators table and its indexes

-- Drop trigger
DROP TRIGGER IF EXISTS update_business_trip_verificators_updated_at ON business_trip_verificators;

-- Drop table
DROP TABLE IF EXISTS business_trip_verificators;