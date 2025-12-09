-- Rollback migration: Drop business trip histories table

-- Drop indexes
DROP INDEX IF EXISTS idx_business_trip_histories_changed_by_user_id;
DROP INDEX IF EXISTS idx_business_trip_histories_created_at;
DROP INDEX IF EXISTS idx_business_trip_histories_change_type;
DROP INDEX IF EXISTS idx_business_trip_histories_business_trip_id;

-- Drop table
DROP TABLE IF EXISTS business_trip_histories;
